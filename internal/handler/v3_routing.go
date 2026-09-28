// v3_routing.go 鈥?v3u/v3 op 鍒嗗彂妗嗘灦銆?
//
// op 璺敱鏄犲皠琛ㄤ粠瀹樻柟鍓嶇 bundle (162 鍙傜収鏈轰笅鍙戠殑 assets/index-*.js) 閫嗗悜鐢熸垚,
// 鍏?267 涓敮涓€ op 鐮併€傚畬鏁存竻鍗曡 .temp/_op_routing_map.txt銆?
//
// 鍒嗗眰绛栫暐:
//   绗竴灞? 宸插疄鐜扮殑 op 鈫?鐩存帴鏄犲皠鍒扮幇鏈?handler 璺緞
//   绗簩灞? 鏈疄鐜?op  鈫?杩斿洖鍗犱綅鍝嶅簲 (success:true, data:null), 璁╁墠绔笉宕╂簝
//
// 绗﹀彿绾﹀畾: 陇 = U+00A4 (user/user+admin 娣峰悎), 搂 = U+00A7 (admin 涓轰富)

package handler

import (
	"net/http"
	"net/url"
	"sync"
)

// HTTP method 甯搁噺鍒悕 (缂╃煭鏄犲皠琛ㄤ功鍐?
const (
	v3GET    = http.MethodGet
	v3POST   = http.MethodPost
	v3PUT    = http.MethodPut
	v3DELETE = http.MethodDelete
	v3PATCH  = http.MethodPatch
)

// V3Target 鏄竴涓?op 瀵瑰簲鐨勫唴閮ㄨ矾鐢辩洰鏍囥€?
type V3Target struct {
	Method string
	Path   string
}

// v3UserRoutes: /api/v3u 鐨?op 鈫?璺敱 (陇 鍓嶇紑 op, 鐢ㄦ埛渚?
var v3UserRoutes = map[string]V3Target{
	// --- 鐧诲綍/璁よ瘉 (鐧诲綍鍓嶆祦绋? ---
	"¤d8617d022a414fd0": {v3POST, "/api/login"},
	"¤a31d3f4f81bdbb7e": {v3GET, "/api/setup/status"},
	"¤82bd8d71746f1694": {v3POST, "/api/setup/init"},
	"¤602549901192a110": {v3GET, "/api/captcha/config"},
	"¤e659c6cd1e9b15af": {v3POST, "/api/domain-check"},
	"¤8eff17ea3b98abd6": {v3POST, "/api/login/recovery"},
	"¤932a9149769e9ffb": {v3POST, "/api/login/2fa"},

	// --- 鐢ㄦ埛鏍稿績 ---
	"¤ba8ad99c09e3549e": {v3GET, "/api/user/profile"},
	"¤f03140ec862a3778": {v3GET, "/api/user/config"},
	"¤c297bcbaa22132e9": {v3PUT, "/api/user/config"},
	"¤4e17a2ea1eae2f78": {v3GET, "/api/user/token"},
	"¤7db4904acbba8324": {v3POST, "/api/user/token"},
	"¤76be241e26965f78": {v3POST, "/api/user/token"},
	"¤749915ca02ff68ac": {v3POST, "/api/user/password"},
	"¤add70c9ed64935de": {v3GET, "/api/subscriptions"},
	"¤c1b0ede30e8cb370": {v3GET, "/api/user/package-assignments"},
	"¤f616b30ab24ed987": {v3GET, "/api/user/forward-nodes"},
	"¤ff8804c1e493810f": {v3GET, "/api/user/forward-nodes"},
	"¤a14fe83e4851bb18": {v3GET, "/api/user/routed-outbounds"},
	"¤8730410b8d1f6fc0": {v3GET, "/api/user/2fa/status"},
	"¤ed0ab02c1ed8a767": {v3POST, "/api/user/2fa/enable"},
	"¤c8518db672c88609": {v3POST, "/api/user/2fa/disable"},
	"¤2538e7ac038b08d9": {v3GET, "/api/traffic/summary"},
	"¤2033390a35b9e02e": {v3GET, "/api/user/telegram-binding"},
	"¤847da72b1ee38d10": {v3POST, "/api/user/telegram-binding/create"},
	"¤cb2e36664d12dee9": {v3DELETE, "/api/user/telegram-binding/delete"},
	"¤056367cb172aa187": {v3POST, "/api/user/renewal/submit"},
	"¤3ac79622892fe222": {v3GET, "/api/user/renewal-request"},
	"¤b0152fcd392f3a06": {v3PUT, "/api/user/profile/update"},
	"¤1338289308c8db86": {v3POST, "/api/2fa/verify"},

	// --- 妯℃澘 ---
	"¤147319194b493bb4": {v3GET, "/api/user/default-template"},
	"¤db790da69751971e": {v3PUT, "/api/user/default-template"},

	// --- 鍏紑 ---
	"¤c646fd4b5ac6b72b": {v3GET, "/api/public/refetch-interval"},
	"¤fff0c6677d0e5bd4": {v3GET, "/api/branding"},
	"¤f34f9dd7b926d1ff": {v3GET, "/api/public/license-status"},
	"¤c7a3dd71032e909d": {v3GET, "/api/user/profile"},

	// --- admin ops 娣峰叆 (陇 鍓嶇紑浣嗗疄闄呮槸 admin 绔偣, admin 鐧诲綍鍚庡彲缁?v3u 鍙? ---
	"¤030e76b59a09e459": {v3POST, "/api/admin/api-tokens/create"},
	"¤06fe64b7b96ae168": {v3POST, "/api/admin/external-sync/import"},
	"¤10b317f626a031fc": {v3GET, "/api/admin/api-token"},
	"¤4eb5d58779eb6e34": {v3POST, "/api/admin/routed-outbounds/create"},
	"¤677266bc6d7c2ff8": {v3POST, "/api/user/forward/create"},
	"¤adf4e861f21d5408": {v3POST, "/api/admin/external-sync/select"},
	"¤ae8f44d8b265a6c5": {v3GET, "/api/admin/proxy-provider-configs"},
	"¤b6b143a1593924ba": {v3POST, "/api/admin/proxy-provider-configs/create"},
	"¤a6a184282b27df62": {v3POST, "/api/admin/subscribe-files/preview"},
	"¤ea2cf021140a96eb": {v3GET, "/api/admin/proxy-group-categories"},

	"¤fefbb40263a6db5a": {v3GET, "/api/admin/external-subscriptions"},
}

// v3AdminRoutes: /api/v3 鐨?op 鈫?璺敱 (搂 鍓嶇紑 op)
var v3AdminRoutes = map[string]V3Target{
	// 鐢ㄦ埛绠＄悊 鈥?/api/admin/users GET 鍒楄〃 + POST 鍔ㄤ綔 (payload 鍐?action 瀛楁鍖哄垎,
	// 瑙?users.go / mmwf_fresh 鐨勭幇鏈夊疄鐜?銆傚涓?op 鍏变韩鍚屼竴璺緞, 鍔ㄤ綔鍦?payload 涓€?
	"§03f3e85b3feaae55": {v3POST, "/api/admin/users"},
	"§5a8525db88308ba9": {v3POST, "/api/admin/users"},
	"§66811e391bb8c52d": {v3POST, "/api/admin/users"},
	"§89bf7f4e7afb5b90": {v3POST, "/api/admin/users"},
	"§939ed195d70303a3": {v3POST, "/api/admin/users"},
	"§9827bede7148e683": {v3POST, "/api/admin/users"},
	"§984718ccf66fbae7": {v3POST, "/api/admin/users"},
	"§b5c31fab7c9e92f8": {v3POST, "/api/admin/users"},
	"§bd8e679ab4130ab0": {v3POST, "/api/admin/users"},
	"§c54c289c1f6b8352": {v3POST, "/api/admin/users"},
	"§cfb75e2e437f9ac1": {v3POST, "/api/admin/users"},
	"§ef5c012ed717a2d2": {v3POST, "/api/admin/users"},
	"§c8383d4437b6ddff": {v3POST, "/api/admin/users"},
	"§99cc845bbc10f1ec": {v3PUT, "/api/admin/users"},
	"§b7194dfabc229019": {v3PUT, "/api/admin/users"},
	"§f47eb450071f7043": {v3PUT, "/api/admin/users"},
	"§b6f8c98d4b8f1c70": {v3GET, "/api/admin/users"},
	"§f555801cf00bd100": {v3DELETE, "/api/admin/users"},

	// 鑺傜偣绠＄悊
	"§b02ec184f40f46f3": {v3GET, "/api/admin/nodes"},
	"§87069e81b99750b3": {v3POST, "/api/admin/nodes/normalize"},
	"§914b57fe83b5c5cf": {v3POST, "/api/admin/nodes/import"},
	"§994f8315797ab697": {v3POST, "/api/admin/nodes/batch-delete"},
	"§a9032ee81f812843": {v3POST, "/api/admin/nodes/delete"},
	"§a9c7e29900e7e941": {v3POST, "/api/admin/nodes/rename"},
	"§ed699dfa96d688a1": {v3POST, "/api/admin/nodes/update"},
	"§25e0cf8cacb34efb": {v3POST, "/api/admin/nodes/probe"},
	"§69e57ddaa8d7dc42": {v3POST, "/api/admin/nodes/batch-probe"},
	"§0e49bf03631c89fb": {v3POST, "/api/admin/nodes/generate-link"},
	"§1e5a52a405eb9f25": {v3POST, "/api/admin/nodes/fetch-url"},
	"§4fbde412961e348b": {v3POST, "/api/admin/nodes/parse-content"},

	// 杩滅▼鏈嶅姟鍣?
	"§b16e74baa40c6a77": {v3GET, "/api/admin/remote-servers"},
	"§ecc3464be116181c": {v3PUT, "/api/admin/remote-servers/update"},
	"§30ce8ad43219ce97": {v3POST, "/api/admin/remote-servers/delete"},

	// 璁㈤槄鏂囦欢
	"§cce38361e0ab6820": {v3GET, "/api/admin/subscribe-files"},
	"§28b2db0b3d098bc6": {v3POST, "/api/admin/subscribe-files"},
	"§270ee7b84b86817d": {v3PUT, "/api/admin/subscribe-files"},
	"§25d490345d5b787a": {v3DELETE, "/api/admin/certificates"},

	// 璁㈤槄妯℃澘
	"§b0c68e7bd88510c3": {v3GET, "/api/admin/templates"},
	"§d7aade6b83f868d5": {v3POST, "/api/admin/templates/apply"},
	"§96f20dab326c5f02": {v3POST, "/api/admin/templates/preview"},
	"§47b927f4f57ddc14": {v3POST, "/api/admin/templates/parse"},
	"§c1f999de86f21f35": {v3POST, "/api/admin/templates/rename"},
	"§5002f68653e44eec": {v3POST, "/api/admin/templates/save"},
	"§1732a9645baf2104": {v3PUT, "/api/admin/default-template"},
	"§df4763ad40808e15": {v3GET, "/api/admin/default-template"},
	"§114ce6bf92480cd0": {v3PUT, "/api/admin/rule-templates/visibility"},

	// 绯荤粺璁剧疆
	"§1d0cdb7046335717": {v3GET, "/api/admin/branding"},
	"§34d865d8754bf54a": {v3POST, "/api/admin/branding"},
	"§1d7335cef7bbfbbe": {v3PUT, "/api/admin/update-cdn-enabled"},
	"§36ad798e97b8eb75": {v3GET, "/api/admin/update-cdn-enabled"},
	"§6ad59db52d1a7555": {v3GET, "/api/admin/tgbot-settings"},
	"§6fc4d93f1494ff17": {v3PUT, "/api/admin/tgbot-settings/update"},
	"§af7bfc337ce5d5f0": {v3PUT, "/api/admin/default-theme-settings"},
	"§b8dc91f573b49bce": {v3GET, "/api/admin/default-theme-settings"},
	"§1513a9c4775e7503": {v3PUT, "/api/admin/redeem-template"},
	"§18fb86405e2ac501": {v3GET, "/api/admin/redeem-template"},
	"§2201c75d3d30989d": {v3PUT, "/api/admin/system-intervals"},
	"§e27c1db7d09c6552": {v3GET, "/api/admin/system-intervals"},
	"§19708c2f6572f113": {v3PUT, "/api/admin/user-permissions-config/update"},
	"§e8e4f79ba0371c33": {v3GET, "/api/admin/user-permissions-config"},
	"§12483c1b0adaaca3": {v3PUT, "/api/admin/miaomiaowu-features/update"},
	"§cacc77d224f60fca": {v3GET, "/api/admin/miaomiaowu-features"},
	"§c85aaaae883b16c0": {v3PUT, "/api/admin/require-encryption"},
	"§6e8a40e20eb15fb3": {v3GET, "/api/admin/require-encryption"},
	"§70c8a6048196fa74": {v3GET, "/api/admin/silent-mode"},
	"§eed187fcadec562e": {v3PUT, "/api/admin/silent-mode"},

	// --- admin 浠〃鐩?/ 棣栭〉鍔犺浇 ---
	"§0fbcc9f59f37fb01": {v3GET, "/api/admin/traffic-snapshots"},
	"§a1fe0c82e099e01a": {v3GET, "/api/admin/traffic-snapshots"},
	"§1ee6e073a783477a": {v3GET, "/api/admin/server-period-totals"},
	"§4d19022efe32112e": {v3GET, "/api/admin/node-totals"},
	"§bdc12e56331cac06": {v3GET, "/api/admin/user-period-ledger"},
	"§30c661ab32b496fb": {v3GET, "/api/admin/override-scripts"},
	"§f6d02db56338248d": {v3GET, "/api/admin/update/check"},
	"§b7968c44cbade85c": {v3GET, "/api/admin/remote-servers"},
	"§f266a67c244b8e4b": {v3GET, "/api/admin/user-connections"},
	"§99f83512a2e16456": {v3GET, "/api/admin/user-connections"},

	// --- REALITY 鍩熷悕鍏变韩 ---
	"§e5620d7808937bac": {v3GET, "/api/admin/reality-share/status"},
	"§bee6aad5bc78875b": {v3POST, "/api/admin/reality-share/toggle"},
	"§0c8151f70733e471": {v3GET, "/api/admin/dns-providers"},
	"§0c8d3cd4401c5e5f": {v3GET, "/api/admin/certificates"},
	"§10d83c386f8d1b7c": {v3POST, "/api/admin/dns-providers/create"},
	"§10fc4ada7df11915": {v3POST, "/api/admin/certificates/renew"},
	"§425e7f006c69171d": {v3POST, "/api/admin/certificates/deploy-custom"},
	"§7dabbed23841c992": {v3POST, "/api/admin/certificates/auto-deploy"},
	"§9b0b9c74197ba21f": {v3POST, "/api/admin/certificates/deploy"},
	"§9e26dd8266a1b0a1": {v3POST, "/api/admin/certificates/revoke"},
	"§a957dbeb9cfcbe31": {v3POST, "/api/admin/certificates/apply"},
	"§b5880d05be90fa29": {v3POST, "/api/admin/certificates/deploy-master"},
	"§de919472f891dc55": {v3GET, "/api/admin/certificates/master-status"},
	"§ec58e5066090ffcd": {v3POST, "/api/admin/certificates/auto-renew"},
	"§4cb3efd640ac3fda": {v3GET, "/api/admin/certificates/for-forward"},
	"§4eef1855789ba971": {v3GET, "/api/admin/external-https"},
	"§0f8be90ac1e3de32": {v3PUT, "/api/admin/external-https"},
	"§4f32c47c302dc217": {v3POST, "/api/admin/https/enable"},
	"§252bb08625ec2ae7": {v3GET, "/api/admin/security-settings"},
	"§4b029e76c94a1ace": {v3PUT, "/api/admin/security-settings"},
	"§9fa3ac00dfe4f0fc": {v3PUT, "/api/admin/subscription-output-format"},
	"§daabe042dcfe9ceb": {v3GET, "/api/admin/subscription-output-format"},
	"§3b089f71f07ea460": {v3PUT, "/api/admin/probe-disguise-settings"},
	"§4faacba09e0608d4": {v3GET, "/api/admin/probe-disguise-settings"},
	"§3e1656c0252d664f": {v3PUT, "/api/admin/login-wallpaper"},
	"§7d241999ca0f4f2c": {v3GET, "/api/admin/login-wallpaper"},
	"§2f2cf00f36ad1229": {v3PUT, "/api/admin/override-scripts/enabled"},
	"§38e3e2dbc289d507": {v3GET, "/api/admin/short-link-enabled"},
	"§b87e73cbacf521fc": {v3PUT, "/api/admin/short-link-enabled"},
	"§d3669601beb06535": {v3PUT, "/api/admin/node-name-multiplier-prefix"},
	"§d6275f1c7a4f96b3": {v3GET, "/api/admin/node-name-multiplier-prefix"},
	"§ce88c07f2d70cc5c": {v3GET, "/api/admin/probe-cdn-regions"},
	"§deb902cfa363a093": {v3GET, "/api/admin/announcement-config"},
	"§2822011fe090aa61": {v3PUT, "/api/admin/announcement-config"},
	"§24b084abc077d365": {v3GET, "/api/admin/announcements/active"},
	"§40857fe46aeb05ff": {v3POST, "/api/admin/announcements/create"},
	"§ab728f92488cfa78": {v3GET, "/api/admin/notify-config"},
	"§73dde0e9bd041cc6": {v3PUT, "/api/admin/notify-config"},
	"§a460ebe3e259edf7": {v3POST, "/api/admin/notify/template"},
	"§2b9834743dc30110": {v3POST, "/api/admin/telegram/test"},
	"§b8d1fd1558256633": {v3GET, "/api/admin/api-token"},
	"§78ef0a27b9d3f6fc": {v3POST, "/api/admin/api-token/regenerate"},
	"§f5e79256e1217821": {v3GET, "/api/admin/database-status"},
	"§2e2b7627d17ccd9a": {v3GET, "/api/admin/database-migration/progress"},
	"§03d1a7cb4554ffff": {v3POST, "/api/admin/database-migrate"},
	"§bb23cddb0a92ba41": {v3POST, "/api/admin/database-migrate/execute"},
	"§3a986427a9504b14": {v3POST, "/api/admin/master-recovery"},
	"§2e0f1e97b5cd72e5": {v3PUT, "/api/admin/master-url"},
	"§6beeb27b8bed4c78": {v3GET, "/api/admin/master-url"},
	"§b27a3bfd3394f32e": {v3POST, "/api/admin/turnstile/verify"},
	"§c6778fc77f035fbf": {v3GET, "/api/admin/license/settings"},
	"§dca925941b97fca8": {v3GET, "/api/admin/license"},
	"§eed66bc4d6ef6ae6": {v3PUT, "/api/admin/license/key"},
	"§c14411003bbfd9c2": {v3GET, "/api/admin/license/usage"},
	"§6db0ee920bb820d8": {v3GET, "/api/admin/license-badge/display"},
	"§a76b9a8469762779": {v3PUT, "/api/admin/license-badge/display"},
	"§7723a1cbc482408e": {v3POST, "/api/admin/reality-share/withdraw"},
	"§06b70d94d6b6d688": {v3GET, "/api/admin/node-tags"},
	"§bfe785923db4dd0b": {v3GET, "/api/admin/node-uris"},
	"§4f398f72fed71f4a": {v3GET, "/api/admin/nodes/blocked"},
	"§4e49e7f5cb8ed92b": {v3POST, "/api/admin/nodes/chain-proxy"},
	"§6545dddef2f67023": {v3POST, "/api/admin/nodes/chain-proxy/create"},
	"§4b880e4b2834cee2": {v3POST, "/api/admin/speed-test/run-all"},
	"§6827165a12076eda": {v3GET, "/api/admin/speed-testers/update-info"},
	"§df22163cf1eab59d": {v3GET, "/api/admin/speed-testers"},
	"§719eb9f9803571aa": {v3POST, "/api/admin/speed-testers/create"},
	"§d2a8f72e911d9fbb": {v3POST, "/api/admin/speed-testers/revoke"},
	"§d9231fac687ed781": {v3POST, "/api/admin/speed-testers/reset-token"},
	"§f1fd55857a65d6aa": {v3POST, "/api/admin/speedtest/start"},
	"§5403591a02d5bf2e": {v3PUT, "/api/admin/package-node-traffic-name"},
	"§eef2b509864a6b3d": {v3GET, "/api/admin/package-node-traffic-name"},
	"§b41025b6b41d7ede": {v3PUT, "/api/admin/nodes/rename-inline"},
	"§cc8397939ab5c99b": {v3POST, "/api/admin/nodes/sync-address"},
	"§ccd571784ac773ea": {v3GET, "/api/admin/tunnels/routing-resolve"},
	"§d836476040c7676d": {v3POST, "/api/admin/tunnels/create"},
	"§b50bed5f07b14cf0": {v3GET, "/api/admin/node-probe"},
	"§c697b1676f4690a6": {v3POST, "/api/admin/node-probe/save"},
	"§f406cdaea1ab4c5d": {v3PUT, "/api/admin/node-probe/update"},
	"§482d4466ef902059": {v3GET, "/api/admin/xray-servers"},
	"§79961874fd4b5b7c": {v3GET, "/api/admin/xray-outbounds"},
	"§7c5e52451d5ad506": {v3GET, "/api/admin/remote-servers/nics"},
	"§422462a286b38c57": {v3POST, "/api/admin/remote-servers/delete"},
	"§41c6b338fe4e20c2": {v3POST, "/api/admin/server-share/create"},
	"§4f4c710dd859ffb1": {v3POST, "/api/admin/server-share/accept"},
	"§947a588916ea943f": {v3POST, "/api/admin/xray-routing/add-rule"},
	"§c5a8e89653aeceea": {v3PUT, "/api/admin/traffic-stats-servers"},
	"§74e2f85fcc153a24": {v3GET, "/api/admin/agent-websites"},
	"§1f2c28e030897639": {v3POST, "/api/admin/agent-websites/create"},
	"§33754abf5250bea2": {v3POST, "/api/admin/agent-websites/detect"},
	"§cc393d02831158f0": {v3DELETE, "/api/admin/agent-websites/delete"},
	"§d927e7f3010e60c9": {v3GET, "/api/admin/forward/chains"},
	"§ce24c9cd335994fe": {v3GET, "/api/admin/forward/groups"},
	"§6a7a41a553bc180a": {v3POST, "/api/admin/forward/chains/create"},
	"§c4d4cc7f618ee06e": {v3POST, "/api/admin/forward/backends/create"},
	"§5f2a6b70ac6ce6c9": {v3POST, "/api/admin/forward/test-connection"},
	"§159cccaa03f52855": {v3POST, "/api/admin/inbound-wizard/test-domain"},
	"§23ddbba0c94f6f94": {v3POST, "/api/admin/inbound-wizard/self-signed-cert"},
	"§98d09c83a1e8d7c8": {v3POST, "/api/admin/inbound-wizard/generate-reality-keys"},
	"§729bf9f3c4d92fe8": {v3POST, "/api/admin/inbound-wizard/generate-encryption"},
	"§930b2c9494c027d6": {v3POST, "/api/admin/inbound-wizard/update-server"},
	"§ae688dc8afd80d01": {v3POST, "/api/admin/inbound-wizard/remove-domain"},
	"§f3c61351ebd80ad0": {v3POST, "/api/admin/inbound-wizard/restore-domain"},
	"§273f84b3601ec985": {v3POST, "/api/admin/users/update-email"},
	"§97e3e31737510104": {v3GET, "/api/admin/packages"},
	"§e917e65c1964a8e3": {v3POST, "/api/admin/packages/create"},
	"§41cc1a7b5ae7df46": {v3POST, "/api/admin/packages/update"},
	"§8449163f7a3871b3": {v3POST, "/api/admin/packages/create-price"},
	"§8e84cec95bcf3f7e": {v3DELETE, "/api/admin/packages/cancel"},
	"§67cf99d2a99bc245": {v3GET, "/api/admin/users/subscription-packages"},
	"§58c49471ca05b62e": {v3GET, "/api/admin/tg-bot-invites"},
	"§e410d90f2ef31f09": {v3POST, "/api/admin/tg-bot-invites/create"},
	"§7a6d317058456dcd": {v3POST, "/api/admin/tg-bot-invites/revoke"},
	"§08ff7bfb146668ba": {v3PUT, "/api/admin/agent-log-enabled"},
	"§f59e4baed1d076fc": {v3GET, "/api/admin/agent-log-enabled"},
	"§9ad03986e56d8b74": {v3GET, "/api/admin/log-files"},
	"§66f5fbe14ac11482": {v3DELETE, "/api/admin/log-files/purge"},
	"§2b48a7e84a7d647e": {v3GET, "/api/admin/agent-log-files"},
	"§51151595c817237d": {v3DELETE, "/api/admin/agent-log-files/purge"},
	"§13d5dff7fc630a5f": {v3GET, "/api/admin/security/whitelist"},
	"§5bf431e1dd96a994": {v3GET, "/api/admin/security/bans"},
	"§a0acbfbc936739c2": {v3POST, "/api/admin/security/ban"},
	"§5e510c6e6579c9fe": {v3PUT, "/api/admin/security/whitelist"},
	"§6e8413e5be417c8c": {v3GET, "/api/admin/task-types"},
	"§7eade1aa7f22e98a": {v3POST, "/api/admin/tasks/run"},
	"§059faf1b9c8c43c4": {v3GET, "/api/admin/rule-files"},
	"§88bb9e88213c4b6f": {v3GET, "/api/admin/custom-rules"},
	"§9a26d474593b59b7": {v3POST, "/api/admin/custom-rules/create"},
	"§e643024c7c39a9fd": {v3GET, "/api/admin/override-scripts"},
	"§7f306a1845f645bd": {v3POST, "/api/admin/override-scripts/create"},
	"§04d85f36b05c8122": {v3GET, "/api/admin/routing-rule-presets"},
	"§6e03775ba638f4b1": {v3POST, "/api/admin/routing-rule-presets/create"},
	"§516a80273ca7f4c1": {v3DELETE, "/api/admin/routing-rule-presets/delete"},
	"§69610e86c8a2642d": {v3GET, "/api/admin/subscribe-files/traffic"},
	"§0dab7283e0a9d58a": {v3GET, "/api/admin/rule-templates"},
	"§25493fc5941face7": {v3GET, "/api/admin/template-v3"},
	"§408bd92a3900a839": {v3GET, "/api/admin/rule-templates/visibility"},
	"§4677fbda90c3c045": {v3POST, "/api/admin/generator/generate"},
	"§48f55c7cb34403e8": {v3POST, "/api/admin/generator/save"},
	"§5cbeed64791d9e8d": {v3POST, "/api/admin/generator/fetch-rule-source"},
	"§d651153042844987": {v3POST, "/api/admin/generator/import-template"},
	"§ee34b2649da3b0bf": {v3POST, "/api/admin/generator/apply-custom-rules"},
	"§1aa598650c6a4049": {v3POST, "/api/admin/proxy-group-categories/sync"},
	"§159c4b1f365db978": {v3POST, "/api/admin/migrate/takeover"},
	"§1ae7e2b6b3f5beb0": {v3POST, "/api/admin/migrate/detect"},
	"§1ee5c456b90eed0e": {v3POST, "/api/admin/migrate/import-node"},
	"§76021b921dbcc0c7": {v3POST, "/api/admin/migrate/import-db"},
	"§b6cea12010d086ad": {v3GET, "/api/admin/migrate/servers"},
	"§ecb9bf0bf12e1ac3": {v3POST, "/api/admin/migrate/online"},
}

// 鍒嗗彂鍣ㄦ寕杞界偣 鈥?鐢?main.go 娉ㄥ叆鍐呴儴 mux (鍚?audit/silent-mode 涓棿浠堕摼),
// 閬垮厤 handler 鍖呭 main 鐨勫惊鐜緷璧栥€?
var (
	v3DispatchMu      sync.RWMutex
	v3DispatchHandler http.Handler
)

// SetV3Dispatcher 瀹夎鍐呴儴璺敱鍒嗗彂鍣ㄣ€?
func SetV3Dispatcher(h http.Handler) {
	v3DispatchMu.Lock()
	v3DispatchHandler = h
	v3DispatchMu.Unlock()
}

// ResolveV3Op 鏌?op 璺敱鏄犲皠琛ㄣ€俰sAdmin=true 璧?admin 琛? false 璧?user 琛ㄣ€?
func ResolveV3Op(op string, isAdmin bool) (V3Target, bool) {
	if isAdmin {
		t, ok := v3AdminRoutes[op]
		return t, ok
	}
	t, ok := v3UserRoutes[op]
	return t, ok
}

// V3Dispatch 鎶婂唴閮ㄨ姹傚杺缁欏畬鏁?mux銆?
func V3Dispatch(w http.ResponseWriter, r *http.Request) {
	v3DispatchMu.RLock()
	h := v3DispatchHandler
	v3DispatchMu.RUnlock()
	if h == nil {
		http.Error(w, "v3 dispatcher not installed", http.StatusInternalServerError)
		return
	}
	h.ServeHTTP(w, r)
}

// parseURLForV3 瑙ｆ瀽璺敱璺緞 (鍚?query)銆?
func parseURLForV3(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		return &url.URL{Path: "/"}
	}
	return u
}

