1. 需求描述

   ccm, csi, vpc-route-controller等组件(后面统称控制器)都需要使用用户的aksk来调用openapi。

   1.1 aksk使用环境变量、secret或者configmap存储，方便配置和更新。 
   1.2 存储的aksk，要支持明文和加密两种方式。
   1.3 存储加密的aksk时，控制器需要秘钥。需要支持在编译控制器时来指定和注入秘钥到控制器二进制，不要把秘钥写死在控制器代码中。
   1.4 ccm, csi, vpc-route-controller不要重复开发以上功能，尽量复用代码。

   基于以上几点需求，我们另开发了aksk-provider通用库，各个控制器通过它来获取有效的aksk。

2. 使用方法

   对于各个控制器来讲，支持两种存放aksk的方式，一种环境变量，一种是文件挂载（包括secret和configmap）。

2.1 环境变量方式

   AkskProvider = env.NewEnvAKSKProvider(encrypt, DefaultCipherKey)
 
   调用env包的NewEnvAKSKProvider函数，需要传入的参数：
   Encrypt：bool类型，是否对sk和securityToken加密了；
   DefaultCipherKey：密钥key（若Encrypt为false，此值为空字符串）

2.2 文件挂载方式

   AkskProvider = file.NewFileAKSKProvider(AkskFilePath, DefaultCipherKey)

   AkskFilePath: 挂载路径，如果为空，则使用默认路径/var/lib/aksk
   DefaultCipherKey：密钥key（若非加密，此值为空字符串）

3. 调用方式

   AkskProvider提供了两个接口，用户获取和更新aksk：

   1) GetAKSK()：获取aksk，aksk的格式见第6部分
   2) ReloadAKSK：重新加载aksk，并获取最新的aksk
