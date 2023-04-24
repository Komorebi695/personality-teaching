package logic

import (
	"fmt"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ms/v20180408"
)

func GetTencentCloudCOSTemporaryCredentials() (string, error) {
	credential := common.NewCredential(
		"AKIDpUwCFD0ZtwwrZ2GG1mMm8UsQ3dyCF4h5",
		"FF0VpuUxEJF1Yq00pVN9jBi182h7JMBk",
	)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "ms.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := ms.NewClient(credential, "", cpf)

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := ms.NewCreateCosSecKeyInstanceRequest()
	response, err := client.CreateCosSecKeyInstance(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return "", nil
	}
	if err != nil {
		panic(any(err))
	}
	// 输出json格式的字符串回包
	fmt.Printf("%s", response.ToJsonString())

	return response.ToJsonString(), nil
}
