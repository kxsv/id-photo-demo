using System;
using System.Collections.Generic;
using System.IO;
using System.Net.Http;
using System.Text.Json;
using System.Threading.Tasks;

namespace IdPhotoDemo;

/// <summary>
/// 鹧应证件照 API 公共客户端。
///
/// 所有接口示例共用的请求封装：
/// - PostAsync()         以 POST 表单方式调用接口
/// - IsSuccess()         判断业务是否成功（code 为 0）
/// - ImageToBase64Async() 读取本地图片并转 Base64
/// - DownloadImageAsync() 下载接口返回的结果图到本地
/// </summary>
public static class ZjzApi
{
    /// <summary>API 网关地址</summary>
    public const string BaseUrl = "https://api.zjzapi.com";

    private const int MaxImageBytes = 8 * 1024 * 1024;

    private static readonly HttpClient Http = new() { Timeout = TimeSpan.FromSeconds(60) };

    private static readonly JsonSerializerOptions JsonOptions = new() { PropertyNameCaseInsensitive = true };

    /// <summary>
    /// 通用返回结构。
    /// </summary>
    public sealed class Response
    {
        public int Code { get; set; }

        public string Msg { get; set; } = string.Empty;

        public JsonElement Data { get; set; }
    }

    /// <summary>
    /// 调用接口。
    /// </summary>
    public static async Task<Response> PostAsync(string path, IDictionary<string, string> form)
    {
        using var content = new FormUrlEncodedContent(form);
        using var response = await Http.PostAsync($"{BaseUrl}/{path.TrimStart('/')}", content);
        var body = await response.Content.ReadAsStringAsync();
        if (!response.IsSuccessStatusCode)
        {
            throw new HttpRequestException($"HTTP {(int)response.StatusCode}：{body}");
        }
        try
        {
            return JsonSerializer.Deserialize<Response>(body, JsonOptions)
                ?? throw new InvalidOperationException($"返回内容不是合法 JSON：{body}");
        }
        catch (JsonException exception)
        {
            throw new InvalidOperationException($"返回内容不是合法 JSON：{body}", exception);
        }
    }

    /// <summary>
    /// 业务是否成功（code 为 0）。
    /// </summary>
    public static bool IsSuccess(Response response) => response.Code == 0;

    /// <summary>
    /// 读取本地图片并转 Base64（不带 data:image/... 前缀）。
    /// </summary>
    public static async Task<string> ImageToBase64Async(string path)
    {
        if (!File.Exists(path))
        {
            throw new FileNotFoundException($"图片不存在：{path}");
        }
        var content = await File.ReadAllBytesAsync(path);
        if (content.Length > MaxImageBytes)
        {
            throw new InvalidOperationException($"图片不能超过 8M：{path}");
        }
        return Convert.ToBase64String(content);
    }

    /// <summary>
    /// 下载接口返回的结果图到本地。
    /// </summary>
    public static async Task<string> DownloadImageAsync(string imageUrl, string saveDirectory, string name)
    {
        Directory.CreateDirectory(saveDirectory);
        var file = Path.Combine(saveDirectory, name);
        using var response = await Http.GetAsync(imageUrl);
        if (!response.IsSuccessStatusCode)
        {
            throw new HttpRequestException($"下载失败：HTTP {(int)response.StatusCode}");
        }
        await using var stream = await response.Content.ReadAsStreamAsync();
        await using var output = File.Create(file);
        await stream.CopyToAsync(output);
        return file;
    }

    /// <summary>
    /// 提取 JSON 字段，数字与字符串两种写法都兼容（数字按原始文本返回）。
    /// 接口部分字段（如 pic_id、order_no、手机/规格数值）有时返回数字、有时返回字符串，
    /// 用本方法统一按字符串读取，避免 GetString() 在数值上抛异常。
    /// </summary>
    public static string GetString(JsonElement element, string name)
    {
        if (!element.TryGetProperty(name, out var value))
        {
            return string.Empty;
        }
        return value.ValueKind == JsonValueKind.Number ? value.GetRawText() : value.GetString() ?? string.Empty;
    }

    /// <summary>
    /// 下载一组「名称 -> 图片URL」的结果图到本地目录，并打印每个结果；
    /// 单个图片下载失败不影响其余图片，仅打印错误后继续。
    /// </summary>
    public static async Task SaveImagesAsync(string saveDirectory, JsonElement listValue)
    {
        if (listValue.ValueKind != JsonValueKind.Object)
        {
            Console.WriteLine("返回数据中没有图片列表。");
            return;
        }
        foreach (var entry in listValue.EnumerateObject())
        {
            var imageUrl = entry.Value.GetString();
            if (string.IsNullOrEmpty(imageUrl))
            {
                continue;
            }
            try
            {
                var file = await DownloadImageAsync(imageUrl, saveDirectory, $"{entry.Name}.jpg");
                Console.WriteLine($"{entry.Name}：{file}");
            }
            catch (Exception exception)
            {
                Console.WriteLine($"{entry.Name} 下载失败：{exception.Message}");
            }
        }
    }
}
