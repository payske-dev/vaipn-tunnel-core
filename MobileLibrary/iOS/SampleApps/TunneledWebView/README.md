# iOS Library Sample App: TunneledWebView

## Tunneling WKWebView

***Only available on iOS 17.0+***

This app tunnels WKWebView traffic by proxying requests through the local Vaipn proxies created by [VaipnTunnel](https://github.com/payske-dev/vaipn-tunnel-core/tree/master/MobileLibrary/iOS/VaipnTunnel).
The listening Vaipn proxy ports can be obtained via TunneledAppDelegate delegate callbacks (see `onListeningSocksProxyPort` and `onListeningHttpProxyPort` in `AppDelegate.swift`).

Tunneling WKWebView traffic was made possible by the introduction of https://developer.apple.com/documentation/network/nwparameters/privacycontext/4156642-proxyconfigurations, which is only available on iOS 17.0+. Tunneling WKWebView traffic on iOS < 17.0 is not supported.

## Tunneling UIWebView

Submissions to App Review that use UIWebView will be rejected by Apple, as it has been fully deprecated. Consequently, our sample app that tunnels UIWebView traffic has been archived, but it can still be found at https://github.com/payske-dev/vaipn-tunnel-core/tree/877122a433a4bc4061b93d2e3f8518900af3e6c7/MobileLibrary/iOS/SampleApps/TunneledWebView.

## *\*\* Caveats \*\*\*

### WKWebView

WKWebView proxying has not been thoroughly tested for the shortcomings that its predecessor, UIWebView proxying with NSURLProtocol, exhibited. These shortcomings are outlined below.

### i18n API Leaks Timezone

The Internationalization API (i18n) provides websites, though a JavaScript API, with access to the timezone used by
the user's browser (in this case UIWebView). This does not reveal the precise location of the user, but can be accurate
enough to identify the city in which the user is located.

Like the "Untunneled WebRTC" issue mentioned below, the i18n API cannot be disabled without disabling JavaScript.         

### NSURLProtocol Challenges

***NSURLProtocol is only partially supported by UIWebView (https://bugs.webkit.org/show_bug.cgi?id=138169) and iOS,
meaning that some network requests are made out of process and are consequently untunneled.***

Below we address the exceptions we have encountered, but there may be more.

### Untunneled Media

***In some versions of iOS audio and video are fetched out of process in mediaserverd and therefore are not intercepted 
by NSURLProtocol.***

*In our limited testing iOS 9/10 leak and iOS 11 does not leak.*

#### Workarounds

***It is worth noting that this fix is inexact and may not always work. If one has control over the HTML being rendered and resources being fetched with XHR it is preferable to alter 
the media source URLs directly beforehand instead of relying on the javascript injection trick.***

***This is a description of a workaround used in the [Vaipn Browser iOS app](https://github.com/Psiphon-Inc/endless) and not of what is implemented in TunneledWebView.
TunneledWebView *does NOT* attempt to tunnel all audio/video content in UIWebView. This is only a hack which allows tunneling
audio and video in UIWebView on versions of iOS which fetch audio/video out of process.***

#### Background
In [PsiphonBrowser](https://github.com/Psiphon-Inc/endless) we have implemented a workaround for audio and video being 
fetched out of process.

[VaipnTunnel's](https://github.com/payske-dev/vaipn-tunnel-core/tree/master/MobileLibrary/iOS/VaipnTunnel/VaipnTunnel)
HTTP Proxy also offers a ["URL proxy (reverse proxy)"](https://github.com/payske-dev/vaipn-tunnel-core/blob/631099d086c7c554a590b0cb76766be6dce94ef9/vaipn/httpProxy.go#L45-L70) 
mode that relays requests for HTTP or HTTPS or URLs specified in the proxy request path. 
 
This reverse proxy can be used by constructing a URL such as `http://127.0.0.1:<proxy-port>/tunneled-rewrite/<origin media URL>?m3u8=true`.

When the retrieved resource is detected to be a [M3U8](https://en.wikipedia.org/wiki/M3U#M3U8) playlist a rewriting rule is applied to ensure all the URL entries
are rewritten to use the same reverse proxy. Otherwise it will be returned unmodified.

#### Fix

* Media element URLs are rewritten to use the URL proxy (reverse proxy).
* This is done by [injecting javascript](https://github.com/Psiphon-Inc/endless/blob/b0c33b4bbd917467a849ad8c51a225c2d4dab260/Endless/Resources/injected.js#L379-L408) 
into the HTML [as it is being loaded](https://github.com/Psiphon-Inc/endless/blob/b0c33b4bbd917467a849ad8c51a225c2d4dab260/External/JiveAuthenticatingHTTPProtocol/JAHPAuthenticatingHTTPProtocol.m#L1274-L1280) 
which [rewrites media URLs to use the URL proxy (reverse proxy)](https://github.com/Psiphon-Inc/endless/blob/b0c33b4bbd917467a849ad8c51a225c2d4dab260/Endless/Resources/injected.js#L319-L377).
* If a [CSP](https://en.wikipedia.org/wiki/Content_Security_Policy) 
is found in the header of the response, we need to modify it to allow our injected javascript to run.
  * This is done by [modifying the
CSP](https://github.com/Psiphon-Inc/endless/blob/b0c33b4bbd917467a849ad8c51a225c2d4dab260/External/JiveAuthenticatingHTTPProtocol/JAHPAuthenticatingHTTPProtocol.m#L1184-L1228) 
to include a nonce generated for our injected javascript, which is [included in the script tag](https://github.com/Psiphon-Inc/endless/blob/b0c33b4bbd917467a849ad8c51a225c2d4dab260/External/JiveAuthenticatingHTTPProtocol/JAHPAuthenticatingHTTPProtocol.m#L1276).

*Requests to localhost (`127.0.0.1`) should be [excluded from being proxied](https://github.com/payske-dev/vaipn-tunnel-core/blob/master/MobileLibrary/iOS/SampleApps/TunneledWebView/External/JiveAuthenticatingHTTPProtocol/JAHPAuthenticatingHTTPProtocol.m#L283-L287) so the system does not attempt to proxy loading the rewritten URLs. They will be correctly proxied through VaipnTunnel's reverse proxy.*

### Untunneled OCSP Requests

See "Online Certificate Status Protocol (OCSP) Leaks" in [../../USAGE.md](../../USAGE.md).

### Untunneled WebRTC

WebRTC in UIWebView does not follow NSURLProtocol and cannot be disabled without disabling JavaScript. If not disabled, 
WebRTC will leak the untunneled client IP address and the WebRTC connection may be performed entirely outside of the
tunnel.

One solution would be to use a WebRTC library which allows setting a proxy; or allows all requests to be intercepted, and
subsequently proxied, through NSURLProtocol.

More details can be found in this issue: https://github.com/OnionBrowser/OnionBrowser/issues/117.

## Configuring, Building, Running

The sample app requires some extra files and configuration before building.

### Get the framework.

1. Run `pod install`

### Get the configuration.

1. Contact Vaipn Inc. to obtain configuration values to use in your app. 
   (This is requried to use the Vaipn network.)
2. Make a copy of `TunneledWebView/vaipn-config.json.stub`, 
   removing the `.stub` extension.
3. Edit `vaipn-config.json`. Remove the comments and fill in the values with 
   those received from Vaipn Inc. The `"ClientVersion"` value is up to you.

### Ready!

TunneledWebView should now compile and run.

### Loading different URLs

Just update `urlString = "https://freegeoip.net"` in `onConnected` to load a different URL into `UIWebView` with TunneledWebView.

## License

See the [LICENSE](../LICENSE) file.
