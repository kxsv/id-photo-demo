// 鹧应证件照 API 公共客户端。
//
// 所有接口示例共用的请求封装：
// - post()            以 POST 表单方式调用接口
// - isSuccess()       判断业务是否成功（code 为 0）
// - imageToBase64()   读取本地图片并转 Base64
// - downloadImage()   下载接口返回的结果图到本地
//
// 仅在服务端运行，请勿将 API Key 放进浏览器等客户端代码。

import fs from 'node:fs/promises';
import path from 'node:path';

export const BASE_URL = 'https://api.zjzapi.com';

const REQUEST_TIMEOUT = 60_000;
const MAX_IMAGE_BYTES = 8 * 1024 * 1024;

/**
 * 调用接口，返回解析后的 JSON。
 */
export async function post(path, params) {
  const form = new URLSearchParams(params);
  const response = await fetch(`${BASE_URL}/${path}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    body: form,
    signal: AbortSignal.timeout(REQUEST_TIMEOUT),
  });
  if (!response.ok) {
    throw new Error(`HTTP ${response.status}：${await response.text()}`);
  }
  return response.json();
}

/**
 * 业务是否成功（code 为 0）。
 */
export function isSuccess(result) {
  return Number(result?.code) === 0;
}

/**
 * 读取本地图片并转 Base64（不带 data:image/... 前缀）。
 */
export async function imageToBase64(imagePath) {
  const content = await fs.readFile(imagePath);
  if (content.length > MAX_IMAGE_BYTES) {
    throw new Error(`图片不能超过 8M：${imagePath}`);
  }
  return content.toString('base64');
}

/**
 * 下载接口返回的结果图到本地。
 */
export async function downloadImage(imageUrl, saveDir, name) {
  await fs.mkdir(saveDir, { recursive: true });
  const response = await fetch(imageUrl, { signal: AbortSignal.timeout(REQUEST_TIMEOUT) });
  if (!response.ok) {
    throw new Error(`下载失败：HTTP ${response.status}`);
  }
  const file = path.join(saveDir, name);
  await fs.writeFile(file, Buffer.from(await response.arrayBuffer()));
  return file;
}
