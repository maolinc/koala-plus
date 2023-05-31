export interface LoginRequestData {
  /** admin 或 editor */
  authType: string
  authKey: string
  /** 密码 */
  password: string
}

export type LoginCodeResponseData = ApiResponseData<string>

export type LoginResponseData = ApiResponseData<{ accessToken: string }>

export type UserInfoResponseData = ApiResponseData<{ username: string; roles: string[] }>
