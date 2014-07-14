/**
 * Copyright (c) 2013 Pixonic.
 * All rights reserved.
 */
package com.pixonic.breakpadintergation;

import java.io.File;
import java.security.MessageDigest;
import java.util.HashMap;
import java.util.Map;
import java.util.UUID;

import org.apache.http.HttpResponse;
import org.apache.http.client.HttpClient;
import org.apache.http.client.methods.HttpPost;
import org.cocos2dx.lib.PIJniCommon;
import org.json.JSONObject;

import android.app.Activity;
import android.app.ProgressDialog;
import android.content.DialogInterface;
import android.content.pm.PackageInfo;
import android.net.http.AndroidHttpClient;
import android.os.AsyncTask;
import android.os.Build;
import android.os.Handler;
import android.os.Looper;
import android.util.Log;

import com.pinidea.ios.sxd.PIConfig;
import com.pinidea.ios.sxd.R;
import com.pinidea.ios.sxd.Road2Immortal;
import com.pinidea.ios.sxd.tools.DeviceUtils;
import com.pinidea.ios.sxd.tools.md5.PIMd5;

/**
 *
 */
public class CrashHandler {

	enum UploadType {
		UT_NATIVE,
		UT_LOG,
		UT_JS,
	}

	private static final String TAG = "CrashHandler";
	private static CrashHandler msSingletonInstance;

	private Activity mActivity;
	private String mSubmitUrl;

	private ProgressDialog mSendCrashReportDialog;
	private static String msApplicationName = "SXD Dump";

	private static HashMap<String, String> optionalFilesToSend = null;
	private static JSONObject optionalParameters = null;


	public static CrashHandler getInstance() {
		return msSingletonInstance;
	}

	public static void init() {
		if (msSingletonInstance == null) {
			msSingletonInstance = new CrashHandler();
		}
	}

	private CrashHandler() {
		mActivity = (Activity) Road2Immortal.context;
		mSubmitUrl = PIConfig.crashSubmitUrl;
		if (msApplicationName == null) {
			msApplicationName = mActivity.getApplicationContext()
					.getPackageName();
		}

		nativeInit(mActivity.getFilesDir().getAbsolutePath());
	}

	/**
	 * Sets a name of the application
	 * 
	 * @param appName
	 *            application name
	 */
	public static void setApplicationName(final String appName) {
		assert (appName != null);
		msApplicationName = appName;
	}

	// / Sets additional file with name `name` to send to server with path
	// `file`
	// / File path needs to be absolute
	public static void includeFile(final String name, final String file) {
		if (optionalFilesToSend == null) {
			optionalFilesToSend = new HashMap<String, String>();
		}

		optionalFilesToSend.put(name, file);
	}

	// / Sets additional request data for dump as json object `params`
	public static void includeJsonData(final JSONObject params) {
		optionalParameters = params;
	}

	// / NATIVE IMPLEMENTATION GLUE ///

	private native void nativeInit(String path);

	/**
	 * A signal handler in native code has been triggered. As our last gasp,
	 * launch the crash handler (in its own process), because when we return
	 * from this function the process will soon exit.
	 */
	static public void nativeCrashed(final String dumpFile) {
		if (msSingletonInstance != null) {
			msSingletonInstance.onCrashed(dumpFile);
		}

		final RuntimeException exception = new RuntimeException(
				"crashed here (native trace should follow after the Java trace)");
		exception.printStackTrace();
		throw exception;
	}

	// / CRASH HANDLING PROCESS ///

	private void onCrashed(final String dumpFile) {
		try {
			createUploadPromtAlert(dumpFile);
			synchronized (this) {
				// lock crashed thread
				wait();
			}
		} catch (final Throwable t) {
			Log.e(TAG, "Error.", t);
		}
		Log.i(TAG, "exit");
	}

	private void createUploadPromtAlert(final String dumpFile) {
		new Thread(new Runnable() {
			@Override
			public void run() {
				// create looper
				Looper.prepare();
				createUploadPromtAlertImpl(dumpFile);

				Looper.loop();
			}
		}).start();
	}

	private void createUploadPromtAlertImpl(final String dumpFile) {
		(new SaveCrashReportTask()).execute(dumpFile);
		return;
	}

	private void onCancelDialog(final DialogInterface dialog) {
		dialog.dismiss();
		finish();
	}

	private void finish() {
		synchronized (this) {
			// release crashed thread
			notifyAll();
		}

		new Handler().post(new Runnable() {
			@Override
			public void run() {
				Looper.myLooper().quit();
			}
		});
	}

	private void createSendDialog() {
		mSendCrashReportDialog = new ProgressDialog(mActivity);
		mSendCrashReportDialog.setMax(100);
		mSendCrashReportDialog.setMessage(mActivity
				.getText(R.string.sending_crash_report));
		mSendCrashReportDialog
				.setProgressStyle(ProgressDialog.STYLE_HORIZONTAL);
		mSendCrashReportDialog.setIndeterminate(false);
		mSendCrashReportDialog.setCancelable(false);
		mSendCrashReportDialog.show();
	}

	private void sendCrashReport(final String dumpFile, final String uuid, UploadType type) {
		String type_str ="";
		switch (type) {
					case UT_NATIVE:
						type_str="MD5";
					break;
					case UT_LOG:
						type_str="LOG";
					break;
					case UT_JS:
						type_str="js";
					break;

					default:
					break;
				}
		(new SendCrashReportTask(mSubmitUrl)).execute(dumpFile, uuid, type_str);
	}

	private class SendCrashReportTask extends
			AsyncTask<String, Integer, Boolean> {

		String mSubmitUrl;

		SendCrashReportTask(String submitUrl) {
			mSubmitUrl = submitUrl;
		}

		protected Boolean doInBackground(String... dumpFiles) {
			sendFile(dumpFiles[0], dumpFiles[1], dumpFiles[2]);
			return true;
		}

		protected void onProgressUpdate(Integer... progress) {
		}

		private void sendFile(String dumpFile, String uuid, String type_str) {
			final SendCrashReportTask task = this;
			HttpClient httpclient = null;
			try {

				httpclient = AndroidHttpClient.newInstance("Breakpad Client");
				HttpPost httppost = new HttpPost(mSubmitUrl);

				MultipartHttpEntity httpEntity = new MultipartHttpEntity();

				String md5srcstr = "UUID:" + uuid + "\ndevice:"
						+ DeviceUtils.getDeviceName() + "\nversion:"
						+ PIJniCommon.getJniLocalVersion() + "\nproduct_name:"
						+ msApplicationName + "\n";

				String md5_str = PIMd5.generateMd5(md5srcstr.getBytes());

		
				httpEntity.addValue(type_str, md5_str);
		
				httpEntity.addValue("UUID", uuid);
				httpEntity.addValue("device", DeviceUtils.getDeviceName());
				httpEntity.addValue("version", PIJniCommon.getJniLocalVersion()
						+ "");
				httpEntity.addValue("product_name", msApplicationName);
				// httpEntity.addValue("report_id", dumpFile.replace(".dmp",
				// ""));
				httpEntity.addFile("symbol_file", "report.dmp", new File(
						Road2Immortal.logDir, dumpFile));

				if (optionalParameters != null) {
					httpEntity.addValue("optional",
							optionalParameters.toString());
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

				Log.v(TAG,
						"request complete, code = "
								+ String.valueOf(resp.getStatusLine()
										.getStatusCode()));
			} catch (final Throwable t) {
				Log.e(TAG, "failed to send file", t);
			} finally {
				if (httpclient != null)
					((AndroidHttpClient) httpclient).close();
			}

			if (dumpFile == "sxddump.dmp") {
				File file = new File(Road2Immortal.logDir, dumpFile);
				file.delete();
			}

			synchronized (this) {
				this.notifyAll();
			}
		}
	}

	private void desptroySendDialog() {
		finish();
	}

	private class SaveCrashReportTask extends
			AsyncTask<String, Integer, Boolean> {

		SaveCrashReportTask() {
		}

		protected Boolean doInBackground(String... dumpFiles) {
			saveFile(dumpFiles[0]);
			return true;
		}

		protected void onPostExecute(Boolean result) {
			desptroySendDialog();
		}

		private void saveFile(String dumpFile) {
			File dump = new File(mActivity.getFilesDir().getAbsolutePath()
					+ "/" + dumpFile);
			dump.renameTo(new File(Road2Immortal.logDir, "sxddump.dmp"));
			File logFile = new File(mActivity.getFilesDir().getAbsolutePath()
					+ "/" + "sxdlog.txt");
			logFile.renameTo(new File(Road2Immortal.logDir, "sxdlog.txt"));
			synchronized (this) {
				this.notifyAll();
			}
		}
	}

	public void UploadDumpFile() {
		new Thread(new Runnable() {
			@Override
			public void run() {
				// create looper
				Looper.prepare();
				String uuid = UUID.randomUUID().toString();
				File file = new File(Road2Immortal.logDir, "sxddump.dmp");
				if (file.exists()) {
					sendCrashReport("sxddump.dmp", uuid, UploadType.UT_NATIVE);

					file = new File(Road2Immortal.logDir, "sxdlog.txt");
					if (file.exists()) {
						sendCrashReport("sxdlog.txt", uuid, UploadType.UT_LOG);
					}
				}
				Looper.loop();
			}
		}).start();
	}
}
