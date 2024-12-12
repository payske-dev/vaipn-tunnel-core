/*
Licensed under Creative Commons Zero (CC0).
https://creativecommons.org/publicdomain/zero/1.0/
*/

package ca.vaipn.tunneledwebview;

import android.content.Context;
import android.os.Bundle;
import androidx.appcompat.app.AppCompatActivity;

import android.util.Log;
import android.webkit.WebSettings;
import android.webkit.WebView;
import android.widget.ArrayAdapter;
import android.widget.ListView;

import org.json.JSONException;
import org.json.JSONObject;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.util.List;
import java.util.concurrent.atomic.AtomicInteger;

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
// - The Vaipn tunnel is started in onResume(). VaipnTunnel.start()
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
//   psi-0.0.10.aar into app/libs).
//
// - Use app/src/main/java/ca/vaipn/VaipnTunnel.java, which provides
//   a higher-level wrapper around the Vaipn Library module. This file
//   shows how to use VaipnTunnel and VaipnTunnel.TunneledApp.
//
//----------------------------------------------------------------------------------------------

public class MainActivity extends AppCompatActivity
        implements VaipnTunnel.HostService {

    private static final String TAG = "TunneledWebView";
    private ListView mListView;
    private WebView mWebView;

    private ArrayAdapter<String> mLogMessages;
    private AtomicInteger mLocalHttpProxyPort;

    private VaipnTunnel mVaipnTunnel;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        mListView = (ListView)findViewById(R.id.listView);
        mWebView = (WebView)findViewById(R.id.webView);
        WebSettings webSettings = mWebView.getSettings();
        webSettings.setJavaScriptEnabled(true);

        mLogMessages = new ArrayAdapter<String>(
                this, R.layout.log_message, R.id.logMessageTextView);

        mListView.setAdapter(mLogMessages);

        mLocalHttpProxyPort = new AtomicInteger(0);

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

        // NOTE: functions called via VaipnTunnel.TunneledApp may be
        // called on background threads. It's important to ensure that
        // these threads are not blocked and that UI functions are not
        // called directly from these threads. Here we use runOnUiThread
        // to handle this.

        runOnUiThread(new Runnable() {
            public void run() {
                WebViewProxySettings.setLocalProxy(
                        MainActivity.this, mLocalHttpProxyPort.get());
                mWebView.loadUrl("https://psip.me/");
            }
        });
    }

    private void logMessage(final String message) {
        runOnUiThread(new Runnable() {
            public void run() {
                mLogMessages.add(message);
                mListView.setSelection(mLogMessages.getCount() - 1);
                Log.d(TAG, "logMessage: " + message);
            }
        });
    }

    //----------------------------------------------------------------------------------------------
    // VaipnTunnel.TunneledApp implementation
    //
    // NOTE: these are callbacks from the Vaipn Library
    //----------------------------------------------------------------------------------------------

    @Override
    public Context getContext() {
        return this;
    }

    @Override
    public String getVaipnConfig() {
        try {
            JSONObject config = new JSONObject(
                    readInputStreamToString(
                            getResources().openRawResource(R.raw.vaipn_config)));

            return config.toString();

        } catch (IOException e) {
            logMessage("error loading Vaipn config: " + e.getMessage());
        } catch (JSONException e) {
            logMessage("error loading Vaipn config: " + e.getMessage());
        }
        return "";
    }

    @Override
    public void onDiagnosticMessage(String message) {
        android.util.Log.i(getString(R.string.app_name), message);
        logMessage(message);
    }

    @Override
    public void onAvailableEgressRegions(List<String> regions) {
        for (String region : regions) {
            logMessage("available egress region: " + region);
        }
    }

    @Override
    public void onSocksProxyPortInUse(int port) {
        logMessage("local SOCKS proxy port in use: " + Integer.toString(port));
    }

    @Override
    public void onHttpProxyPortInUse(int port) {
        logMessage("local HTTP proxy port in use: " + Integer.toString(port));
    }

    @Override
    public void onListeningSocksProxyPort(int port) {
        logMessage("local SOCKS proxy listening on port: " + Integer.toString(port));
    }

    @Override
    public void onListeningHttpProxyPort(int port) {
        logMessage("local HTTP proxy listening on port: " + Integer.toString(port));
        setHttpProxyPort(port);
    }

    @Override
    public void onUpstreamProxyError(String message) {
        logMessage("upstream proxy error: " + message);
    }

    @Override
    public void onConnecting() {
        logMessage("connecting...");
    }

    @Override
    public void onConnected() {
        logMessage("connected");
        loadWebView();
    }

    @Override
    public void onConnectedServerRegion(String region) {
        logMessage("connected server region: " + region);
    }

    @Override
    public void onHomepage(String url) {
        logMessage("home page: " + url);
    }

    @Override
    public void onClientUpgradeDownloaded(String filename) {
        logMessage("client upgrade downloaded");
    }

    @Override
    public void onClientIsLatestVersion() {

    }

    @Override
    public void onUntunneledAddress(String address) {
        logMessage("untunneled address: " + address);
    }

    @Override
    public void onBytesTransferred(long sent, long received) {
        logMessage("bytes sent: " + Long.toString(sent));
        logMessage("bytes received: " + Long.toString(received));
    }

    @Override
    public void onStartedWaitingForNetworkConnectivity() {
        logMessage("waiting for network connectivity...");
    }

    @Override
    public void onStoppedWaitingForNetworkConnectivity() {
        logMessage("finished waiting for network connectivity...");
    }

    @Override
    public void onActiveAuthorizationIDs(List<String> authorizations) {

    }

    @Override
    public void onExiting() {

    }

    @Override
    public void onClientRegion(String region) {
        logMessage("client region: " + region);
    }

    @Override
    public void onClientAddress(String address) {
        logMessage("client address: " + address);
    }

    private static String readInputStreamToString(InputStream inputStream) throws IOException {
        return new String(readInputStreamToBytes(inputStream), "UTF-8");
    }

    private static byte[] readInputStreamToBytes(InputStream inputStream) throws IOException {
        ByteArrayOutputStream outputStream = new ByteArrayOutputStream();
        int readCount;
        byte[] buffer = new byte[16384];
        while ((readCount = inputStream.read(buffer, 0, buffer.length)) != -1) {
            outputStream.write(buffer, 0, readCount);
        }
        outputStream.flush();
        inputStream.close();
        return outputStream.toByteArray();
    }
}
