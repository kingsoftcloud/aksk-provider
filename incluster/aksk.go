package incluster

import (
	"fmt"
	"strings"
	"sync"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	prvd "github.com/kingsoftcloud/aksk-provider"
	"github.com/kingsoftcloud/aksk-provider/types"
	"github.com/kingsoftcloud/aksk-provider/utils"
)

const (
	defaultAkSKCMName          = "user-temp-aksk"
	defaultAkSKCMNameSpace     = "kube-system"
	defaultAkSKSecretName      = "user-temp-aksk"
	defaultAkSKSecretNameSpace = "kube-system"
	defaultKubeconfigPath      = "/root/.kube/config"
)

var _ prvd.AKSKProvider = &AKSKProvider{}

type AKSKProvider struct {
	AkskCMName          string
	AkskCMNameSpace     string
	AkskSecretName      string
	AkskSecretNameSpace string
	Clientset           *kubernetes.Clientset
	CipherKey           string
	AkskMap             sync.Map
	AkskSource          string
}

func NewAKSKProviderByKubeConfigFilePath(akskCMName, akskCMNameSpace, akskSecretName, akskSecretNameSpace, cipherKey, kubeconfigPath string) (prvd.AKSKProvider, error) {
	var akskSource string
	if akskCMName != "" && akskCMNameSpace != "" {
		akskSource = "configmap"
	}
	if akskSecretName != "" && akskSecretNameSpace != "" {
		akskSource = "secret"
	}
	if akskCMName != "" && akskCMNameSpace != "" && akskSecretName != "" && akskSecretNameSpace != "" {
		akskSource = ""
	}
	if kubeconfigPath == "" {
		kubeconfigPath = defaultKubeconfigPath
	}
	if akskCMName == "" {
		akskCMName = defaultAkSKCMName
	}
	if akskCMNameSpace == "" {
		akskCMNameSpace = defaultAkSKCMNameSpace
	}
	if akskSecretName == "" {
		akskSecretName = defaultAkSKSecretName
	}
	if akskSecretNameSpace == "" {
		akskSecretNameSpace = defaultAkSKSecretNameSpace
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	provider := &AKSKProvider{
		Clientset:           clientset,
		AkskCMName:          akskCMName,
		AkskCMNameSpace:     akskCMNameSpace,
		AkskSecretName:      akskSecretName,
		AkskSecretNameSpace: akskSecretNameSpace,
		CipherKey:           cipherKey,
		AkskMap:             sync.Map{},
		AkskSource:          akskSource,
	}
	err = provider.loadAksk()
	if err != nil {
		return nil, err
	}
	return provider, nil
}

func NewAKSKProviderByClientset(akskCMName, akskCMNameSpace, akskSecretName, akskSecretNameSpace, cipherKey string, clientset *kubernetes.Clientset) (prvd.AKSKProvider, error) {
	var akskSource string
	if akskCMName != "" && akskCMNameSpace != "" {
		akskSource = "configmap"
	}
	if akskSecretName != "" && akskSecretNameSpace != "" {
		akskSource = "secret"
	}
	if akskCMName != "" && akskCMNameSpace != "" && akskSecretName != "" && akskSecretNameSpace != "" {
		akskSource = ""
	}
	if akskCMName == "" {
		akskCMName = defaultAkSKCMName
	}
	if akskCMNameSpace == "" {
		akskCMNameSpace = defaultAkSKCMNameSpace
	}
	if akskSecretName == "" {
		akskSecretName = defaultAkSKSecretName
	}
	if akskSecretNameSpace == "" {
		akskSecretNameSpace = defaultAkSKSecretNameSpace
	}
	provider := &AKSKProvider{
		Clientset:           clientset,
		AkskCMName:          akskCMName,
		AkskCMNameSpace:     akskCMNameSpace,
		AkskSecretName:      akskSecretName,
		AkskSecretNameSpace: akskSecretNameSpace,
		CipherKey:           cipherKey,
		AkskMap:             sync.Map{},
		AkskSource:          akskSource,
	}
	err := provider.loadAksk()
	if err != nil {
		return nil, err
	}
	return provider, nil
}

func (pvd *AKSKProvider) loadAksk() error {
	var err error
	if pvd.AkskSource == "configmap" {
		err = pvd.loadAkskInConfigMap()
	} else if pvd.AkskSource == "secret" {
		err = pvd.loadAkskInSecret()
	} else if pvd.AkskSource == "" {
		err = pvd.loadAkskInConfigMap()
		if err != nil {
			err = pvd.loadAkskInSecret()
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func (pvd *AKSKProvider) GetAKSK() (*types.AKSK, error) {
	v, ok := pvd.AkskMap.Load("aksk")
	if ok {
		if !utils.IsExpired(v.(*types.AKSK).ExpiredAt) {
			return v.(*types.AKSK), nil
		} else {
			return nil, fmt.Errorf("aksk expired, please retry")
		}
	}
	return nil, fmt.Errorf("aksk not found")
}

func (pvd *AKSKProvider) ReloadAKSK() (*types.AKSK, error) {
	aksk, err := pvd.GetAKSK()
	if err != nil {
		return nil, err
	}
	pvd.AkskMap.Store("aksk", aksk)
	return aksk, nil
}

func (pvd *AKSKProvider) loadAkskInConfigMap() error {
	var aksk *types.AKSK
	var err error
	aksk, err = utils.GetAkskConfigMap(pvd.AkskCMName, pvd.AkskCMNameSpace, pvd.Clientset)
	if err != nil {
		return err
	}
	if aksk.Cipher != "none" && aksk.Cipher != "" {
		aksk.SK, err = utils.DecryptData(aksk.SK, pvd.CipherKey, aksk.Cipher)
		if err != nil {
			return err
		}
	}
	pvd.AkskMap.Store("aksk", aksk)
	go pvd.watchAkskConfigMap(pvd.AkskCMName, pvd.AkskCMNameSpace)
	return nil
}

func (pvd *AKSKProvider) loadAkskInSecret() error {
	var aksk *types.AKSK
	var err error
	aksk, err = utils.GetAkskSecret(pvd.AkskSecretName, pvd.AkskSecretNameSpace, pvd.Clientset)
	if err != nil {
		return err
	}
	if aksk.Cipher != "none" && aksk.Cipher != "" {
		aksk.SK, err = utils.DecryptData(aksk.SK, pvd.CipherKey, aksk.Cipher)
		if err != nil {
			return err
		}
	}
	pvd.AkskMap.Store("aksk", aksk)
	go pvd.watchAkskSecret(pvd.AkskSecretName, pvd.AkskSecretNameSpace)
	return nil
}

func (pvd *AKSKProvider) watchAkskConfigMap(name, namespace string) {
	informerFactory := informers.NewSharedInformerFactoryWithOptions(
		pvd.Clientset,
		10*time.Hour,
		informers.WithNamespace(namespace),
		informers.WithTweakListOptions(func(options *metav1.ListOptions) {
			options.FieldSelector = fields.OneTermEqualSelector("metadata.name", name).String()
		}))
	informer := informerFactory.Core().V1().ConfigMaps()
	informer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldCm := oldObj.(*v1.ConfigMap)
				newCm := newObj.(*v1.ConfigMap)
				if oldCm.ResourceVersion == newCm.ResourceVersion {
					return
				}
				pvd.loadConfigMap(newCm)
			},
		})
	informer.Informer().Run(wait.NeverStop)
}

func (pvd *AKSKProvider) watchAkskSecret(name, namespace string) {
	informerFactory := informers.NewSharedInformerFactoryWithOptions(
		pvd.Clientset,
		10*time.Hour,
		informers.WithNamespace(namespace),
		informers.WithTweakListOptions(func(options *metav1.ListOptions) {
			options.FieldSelector = fields.OneTermEqualSelector("metadata.name", name).String()
		}))
	informer := informerFactory.Core().V1().Secrets()
	informer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldSecret := oldObj.(*v1.Secret)
				newSecret := newObj.(*v1.Secret)
				if oldSecret.ResourceVersion == newSecret.ResourceVersion {
					return
				}
				pvd.loadSecret(newSecret)
			},
		})
	informer.Informer().Run(wait.NeverStop)
}

func (pvd *AKSKProvider) loadConfigMap(cm *v1.ConfigMap) {
	ak := cm.Data["ak"]
	noDecryptSk := cm.Data["sk"]
	securityToken := cm.Data["securityToken"]
	cipher := cm.Data["cipher"]
	decryptedSk := noDecryptSk
	if cipher != "none" && cipher != "" {
		var err error
		decryptedSk, err = utils.DecryptData(noDecryptSk, pvd.CipherKey, cipher)
		if err != nil {
			klog.Errorf("Failed to decrypt SK: %v", err)
			return
		}
	}
	ts, err := time.Parse(utils.TimeLayoutStr, strings.TrimSpace(cm.Data["expired_at"]))
	if err != nil {
		return
	}
	aksk := &types.AKSK{
		AK:            ak,
		SK:            decryptedSk,
		Cipher:        cipher,
		ExpiredAt:     ts,
		SecurityToken: securityToken,
	}
	pvd.AkskMap.Store("aksk", aksk)
	klog.Infof("ak:%s updated", aksk.AK)
}

func (pvd *AKSKProvider) loadSecret(secret *v1.Secret) {
	ak := string(secret.Data["ak"])
	noDecryptSk := string(secret.Data["sk"])
	securityToken := string(secret.Data["securityToken"])
	cipher := string(secret.Data["cipher"])
	decryptedSk := noDecryptSk
	if cipher != "none" && cipher != "" {
		var err error
		decryptedSk, err = utils.DecryptData(noDecryptSk, pvd.CipherKey, cipher)
		if err != nil {
			klog.Errorf("Failed to decrypt SK: %v", err)
			return
		}
	}
	ts, err := time.Parse(utils.TimeLayoutStr, strings.TrimSpace(string(secret.Data["expired_at"])))
	if err != nil {
		return
	}
	aksk := &types.AKSK{
		AK:            ak,
		SK:            decryptedSk,
		Cipher:        cipher,
		ExpiredAt:     ts,
		SecurityToken: securityToken,
	}
	pvd.AkskMap.Store("aksk", aksk)
	klog.Infof("ak:%s updated", aksk.AK)
}
