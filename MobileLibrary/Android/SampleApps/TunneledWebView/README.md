TunneledWebView README
================================================================================

Overview
--------------------------------------------------------------------------------

TunneledWebView is a sample app that demonstrates embedding the Vaipn Library in
an Android app. TunneledWebView proxies a WebView through the Vaipn tunnel.

Caveats
--------------------------------------------------------------------------------

### i18n API Leaks Timezone

The Internationalization API (i18n) provides websites, though a JavaScript API, with access to the timezone used by
the user's browser (in this case WebView). This does not reveal the precise location of the user, but can be accurate
enough to identify the city in which the user is located.

The i18n API cannot be disabled without disabling JavaScript.

### Untunneled WebRTC

WebRTC requests do not use the configured proxy settings of a WebView. JavaScript must be disabled in a WebView to
effectively disable WebRTC. If not disabled, WebRTC will leak the untunneled client IP address and the WebRTC connection
may be performed entirely outside of the tunnel.

One solution would be to use a WebRTC library which allows setting a proxy; or use 
[Mozilla's GeckoView](https://wiki.mozilla.org/Mobile/GeckoView), which is a WebView alternative which allows disabling
WebRTC.

Integration
--------------------------------------------------------------------------------

Uses the [Vaipn Android Library](../../../Android/README.md).

Integration is illustrated in the main activity source file in the sample app. Here are the key parts.

```Java

/*
 * Copyright (c) 2016, Vaipn Inc.
 * All rights reserved.
 */
 
package ca.vaipn.tunneledwebview;

// ...

import ca.vaipn.VaipnTunnel;

//----------------------------------------------------------------------------------------------
// TunneledWebView
//
// This sample app demonstrates tunneling a WebView through the
// Vaipn Library. This app's main activity shows a log of
// events and a WebView that is loaded once Vaipn is connected.
//
// The flow is as follows:
//
// - The Vaipn tunnel is started in onResume(). VaipnTunnel.startTunneling()
//   is an asynchronous call that returns immediately.
//
// - Once Vaipn has selected a local HTTP proxy listening port, the
//   onListeningHttpProxyPort() callback is called. This app records the
//   port to use for tunneling traffic.
//
// - Once Vaipn has established a tunnel, the onConnected() callback
//   is called. This app now loads the WebView, after setting its proxy
//   to point to Vaipn's local HTTP proxy.
//
// To adapt this sample into your own app:
//
// - Embed a Vaipn config file in app/src/main/res/raw/vaipn_config.
//
// - Add the Vaipn Library AAR module as a dependency (see this app's
//   project settings; to build this sample project, you need to drop
//   ca.vaipn.aar into app/libs).
//----------------------------------------------------------------------------------------------

public class MainActivity extends ActionBarActivity
        implements VaipnTunnel.HostService {

// ...

    @Override
    protected void onCreate(Bundle savedInstanceState) {

        // ...

        mVaipnTunnel = VaipnTunnel.newVaipnTunnel(this);
    }

    @Override
    protected void onResume() {
        super.onResume();

        // NOTE: for demonstration purposes, this sample app
        // restarts Vaipn in onPause/onResume. Since it may take some
        // time to connect, it's generally recommended to keep
        // Vaipn running, so start/stop in onCreate/onDestroy or
        // even consider running a background Service.

        try {
            mVaipnTunnel.startTunneling("");
        } catch (VaipnTunnel.Exception e) {
            logMessage("failed to start Vaipn");
        }
    }

    @Override
    protected void onPause() {
        super.onPause();

        // NOTE: stop() can block for a few seconds, so it's generally
        // recommended to run VaipnTunnel.start()/stop() in a background
        // thread and signal the thread appropriately.

        mVaipnTunnel.stop();
    }

    private void setHttpProxyPort(int port) {

        // NOTE: here we record the Vaipn proxy port for subsequent
        // use in tunneling app traffic. In this sample app, we will
        // use WebViewProxySettings.setLocalProxy to tunnel a WebView
        // through Vaipn. By default, the local proxy port is selected
        // dynamically, so it's important to record and use the correct
        // port number.

        mLocalHttpProxyPort.set(port);
    }

    private void loadWebView() {

        // NOTE: functions called via VaipnTunnel.HostService may be
        // called on background threads. It's important to ensure that
        // these threads are not blocked and that UI functions are not
        // called directly from these threads. Here we use runOnUiThread
        // to handle this.

        runOnUiThread(new Runnable() {
            public void run() {
                WebViewProxySettings.setLocalProxy(
                        MainActivity.this, mLocalHttpProxyPort.get());
                mWebView.loadUrl("https://ipinfo.io/");
            }
        });
    }

    // ...

    //----------------------------------------------------------------------------------------------
    // VaipnTunnel.HostService implementation
    //
    // NOTE: these are callbacks from the Vaipn Library
    //----------------------------------------------------------------------------------------------

    // ...

    @Override
    public void onListeningSocksProxyPort(int port) {
        logMessage("local SOCKS proxy listening on port: " + Integer.toString(port));
    }

    // ...

    @Override
    public void onConnected() {
        logMessage("connected");
        loadWebView();
    }

    // ...

}

```

## License

See the [LICENSE](../LICENSE) file.
