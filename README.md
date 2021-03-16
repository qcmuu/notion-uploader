# notion-uploader 

Upload files to Notion.

上传文件到Notion的小工具。

## Example 

```bash
❯ ./notion-uploader -p ... -t ... install.sh        
Local: /.../install.sh
1.03 KiB / 1.03 KiB [--------------------] 100.00% 810 B p/s 2s
syncing blocks....
Download Link: https://www.notion.so/signed/https%3A%2F%2Fs3-us-west-2.amazonaws.com%2Fsecure.notion-static.com%2F...%2Finstall.sh?table=block&id=...&name=install.sh&userId=...&cache=v2
```

## Parameters

```
-p Page ID
```
Required, the ID in the page link. It is recommended to directly use the secondary page of Workspace so that the program can automatically obtain the Workspace ID, otherwise, you need to specify the Space ID through the `-s` parameter.

```
-t token
```
Required, the `www.notion.so -> token_v2` value in the cookie.

```
-s Workspace ID
```
Optional, applicable to non-secondary/nested pages, manually set Workspace ID

---

上传后默认返回一个自动签名链接，私有页面可以在浏览器登录状态下直接点击下载。对于公开页面的文件链接，可以尝试去掉userid使用，但必须保留id和table两项。

## 上传参数

所有参数不带符号，即形如ce6ad860c0864286a4392d6c2e786e8即可。

```
-p Page ID
```
必须，即页面链接中的那个一大长串的ID。建议直接使用Workspace的次级页面作为上传目标以便程序能自动获取当前Workspace ID，否则需要通过 -s 参数指定 Space ID。

```
-t token
```
必须，即cookie中的www.notion.so -> token_v2项。

```
-s Workspace ID
```
非必须，适用于非次级页面/嵌套的情况，手动设定Workspace ID
