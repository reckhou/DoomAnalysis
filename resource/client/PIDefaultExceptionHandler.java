package com.pixonic.breakpadintergation;

import java.io.ByteArrayOutputStream;
import java.io.File;
import java.io.PrintWriter;
import java.io.StringWriter;
import java.io.Writer;
import java.lang.Thread.UncaughtExceptionHandler;
import java.util.HashMap;
import java.util.Map;
import java.util.UUID;

import org.apache.http.HttpResponse;
import org.apache.http.client.HttpClient;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.util.EntityUtils;
import org.cocos2dx.lib.PIJniCommon;
import org.json.JSONObject;

import android.content.Context;
import android.net.http.AndroidHttpClient;
import android.util.Log;

import com.pinidea.ios.sxd.PIConfig;
import com.pinidea.ios.sxd.tools.md5.PIMd5;
import com.pixonic.breakpadintergation.MultipartHttpEntity;

public class PIDefaultExceptionHandler implements UncaughtExceptionHandler {
	// private static final String FLAG = "CrashHandler";
	private Context mContext;
	private static PIDefaultExceptionHandler INSTANCE = new PIDefaultExceptionHandler();
	private Thread.UncaughtExceptionHandler mDefaultHandler;
	private final static String msApplicationName = "SXD Dump";
	private static JSONObject optionalParameters = null;
	private static HashMap<String, String> optionalFilesToSend = null;

	public static PIDefaultExceptionHandler getInstance() {
		return INSTANCE;
	}

	public void init(Context ctx) {
		mContext = ctx;
		mDefaultHandler = Thread.getDefaultUncaughtExceptionHandler();
		Thread.setDefaultUncaughtExceptionHandler(this);
	}

	@Override
	public void uncaughtException(Thread thread, final Throwable ex) {
		new Thread() {
			public void run() {
				Writer writer = new StringWriter();
				PrintWriter pw = new PrintWriter(writer);
				ex.printStackTrace(pw);
				Log.e("exception", writer.toString());
				sendCrash(writer.toString());
			}
		}.start();
	}

	private void sendCrash(String exceptionStr) {
		HttpClient httpclient = null;
		try {

			httpclient = AndroidHttpClient.newInstance("Breakpad Client");
			HttpPost httppost = new HttpPost(PIConfig.crashSubmitUrl);

			MultipartHttpEntity httpEntity = new MultipartHttpEntity();
			String uuid = UUID.randomUUID().toString();

			String md5srcstr = "UUID:" + uuid + "\ndevice:"
					+ DeviceUtils.getDeviceName() + "\nversion:"
					+ PIJniCommon.getJniLocalVersion() + "\nproduct_name:"
					+ msApplicationName + "\n";
			String md5_str = PIMd5.generateMd5(md5srcstr.getBytes());

			// 去除 除了MINDDUMP的其他信息, 保留注释,备用
			httpEntity.addValue("java", md5_str);
			httpEntity.addValue("UUID", uuid);
			httpEntity.addValue("device", DeviceUtils.getDeviceName());
			httpEntity.addValue("version", PIJniCommon.getJniLocalVersion()
					+ "");
			httpEntity.addValue("product_name", msApplicationName);
			httpEntity.addStream(exceptionStr);
			if (optionalParameters != null) {
				httpEntity.addValue("optional", optionalParameters.toString());
			}

			if (optionalFilesToSend != null) {
				for (final Map.Entry<String, String> file : optionalFilesToSend
						.entrySet()) {
					final File f = new File(file.getValue());
					httpEntity.addFile(file.getKey(), f.getName(), f);
				}
			}

			httpEntity.finish();
			httppost.setEntity(httpEntity);

			httppost.setHeader("Connection", "close");

			// Execute HTTP Post Request
			final HttpResponse resp = httpclient.execute(httppost);

			Log.v("send crash",
					"request complete, code = "
							+ String.valueOf(resp.getStatusLine()
									.getStatusCode()) + "http entity:"
							+ EntityUtils.toString(httpEntity)
							+ "response http entity:"
							+ EntityUtils.toString(resp.getEntity()));
		} catch (final Throwable t) {
			Log.e("send crash exception",
					t.getMessage() + ";" + t.getLocalizedMessage());

		} finally {
			if (httpclient != null)
				((AndroidHttpClient) httpclient).close();
			System.exit(0);
		}
	}

}