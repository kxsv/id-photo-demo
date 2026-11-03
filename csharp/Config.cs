namespace IdPhotoDemo;

/// <summary>
/// Demo 全局配置。
/// </summary>
public static class Config
{
    /// <summary>
    /// 应用的 API Key，在官网「用户中心 - 应用管理」创建应用后获取。
    /// </summary>
    public const string ApiKey = "YOUR_APP_KEY";

    /// <summary>
    /// 用户中心的 access_token，在官网「用户中心 - 账户信息」获取。
    /// 仅「账号相关」接口（account-info、account-app）需要，其余接口只用 ApiKey。
    /// </summary>
    public const string AccessToken = "YOUR_ACCESS_TOKEN";
}
