package oauth

func GetOauthUrl( corp_id string, back_url string, state string, scope string) string {


	if corp_id == ""   {
		return ""
	}
	base_url := "https://open.weixin.qq.com/connect/oauth2/authorize?"

	url := base_url+"appid="+corp_id+"&redirect_uri="+back_url+"&response_type=code&scope="+scope+"&state="+state+"#wechat_redirect"

	return url
}