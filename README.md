1. 需求描述

   ccm, csi, vpc-route-controller等组件(后面统称控制器)都需要使用用户的aksk来调用openapi。
   
   1.1 aksk使用环境变量、secret或者configmap存储，方便配置和更新。
   
   1.2 存储的aksk，要支持明文和加密两种方式。
   
   1.3 存储加密的aksk时，控制器需要秘钥。需要支持在编译控制器时来指定和注入秘钥到控制器二进制，不要把秘钥写死在控制器代码中。
   
   1.4 ccm, csi, vpc-route-controller不要重复开发以上功能，尽量复用代码。
   
   基于以上几点需求，我们另开发了aksk-provider通用库，各个控制器通过它来获取有效的aksk。
   
2. 使用方法
   
   对于各个控制器来讲，支持三种存放aksk的方式，一种环境变量，一种是文件挂载（包括secret和configmap）,最后一种是直接从集群内获取（包括secret或configmap）。
   
   2.1 环境变量方式

   AkskProvider = env.NewEnvAKSKProvider(encrypt, DefaultCipherKey)
   
   调用env包的NewEnvAKSKProvider函数，需要传入的参数：
   
   Encrypt：bool类型，是否对sk加密了；
   
   DefaultCipherKey：密钥key（若Encrypt为false，此值为空字符串）
   
   2.2 文件挂载方式

   AkskProvider = file.NewFileAKSKProvider(AkskFilePath, DefaultCipherKey)
   
   AkskFilePath: 挂载路径，如果为空，则使用默认路径/var/lib/aksk
   
   DefaultCipherKey：密钥key（若非加密，此值为空字符串）

   2.3 inCluster方式

   2.3.1 NewInClusterAKSKProviderByKubeConfigFilePath
   
   AkskProvider = incluster.NewAKSKProviderByKubeConfigFilePath(akskCMName, akskCMNameSpace, akskSecretName, akskSecretNameSpace, cipherKey, kubeconfigPath)

   初始化NewAKSKProviderByKubeConfigFilePath函数，需要传入的参数：
   akskCMName：保存aksk的ConfigMap的name
   akskCMNameSpace：保存aksk的ConfigMap的nameSpace
   akskSecretName：保存aksk的Secret的name
   akskSecretNameSpace：保存aksk的Secret的nameSpace
   cipherKey：密钥key
   kubeconfigPath：kubeconfig的文件路径
   
   2.3.2  NewInClusterAKSKProviderByClientset
   
   AkskProvider = incluster.NewAKSKProviderByClientset(akskCMName, akskCMNameSpace, akskSecretName, akskSecretNameSpace, cipherKey string, clientset)
   
   初始化NewAKSKProviderByClientset函数，需要传入的参数：
   akskCMName：保存aksk的ConfigMap的name
   akskCMNameSpace：保存aksk的ConfigMap的nameSpace
   akskSecretName：保存aksk的Secret的name
   akskSecretNameSpace：保存aksk的Secret的nameSpace
   cipherKey：密钥key
   clientset：k8s集群客户端
   
3. 调用方式
   
    AkskProvider提供了两个接口，用户获取和更新aksk：
   
    GetAKSK()：获取aksk，aksk的格式见第6部分

    ReloadAKSK：重新加载aksk，并获取最新的aksk
