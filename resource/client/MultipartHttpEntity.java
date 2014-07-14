// (C) Pixonic, 2013

package com.pixonic.breakpadintergation;

import java.io.ByteArrayInputStream;
import java.io.File;
import java.io.FileInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.io.SequenceInputStream;
import java.util.ArrayList;
import java.util.Collections;
import java.util.UUID;

import org.apache.http.Header;
import org.apache.http.HttpEntity;
import org.apache.http.message.BasicHeader;
import org.apache.http.protocol.HTTP;

import android.util.Log;

/**
 * Allows you to send in the body of post-request form with sending files More
 * files can attach
 */
public class MultipartHttpEntity implements HttpEntity {
	public static interface ProgressCallback {
		void onProgress(long current, long target);
	}

	private final String BOUNDARY_TAG;

	private static final int BUFFER_SIZE = 2048;
	private static final int EOF_MARK = -1;

	private final ArrayList<InputStream> mInputChuncks = new ArrayList<InputStream>(
			5);
	private long mTotalLength = 0;
	private boolean mReady = false;

	private ProgressCallback progressCallback = null;

	public MultipartHttpEntity() {
		BOUNDARY_TAG = UUID.randomUUID().toString();
	}

	MultipartHttpEntity(ProgressCallback progressCallback) {
		BOUNDARY_TAG = UUID.randomUUID().toString();
		this.progressCallback = progressCallback;
	}

	/**
	 * Adds a string value with given name to this entity
	 * 
	 * @param name
	 *            a name of item
	 * @param value
	 *            a value of item
	 */
	public void addValue(final String name, final String value) {
		final StringBuilder stringBuilder = createHeaderBuilder(name);
		stringBuilder.append(":").append(value).append("\n");

		final String data = stringBuilder.toString();
		mTotalLength += data.length();
		mInputChuncks.add(new ByteArrayInputStream(data.getBytes()));
	}

	/**
	 * Adds a file with given name to this entity
	 * 
	 * @param name
	 *            a name of item
	 * @param fileName
	 *            a name of file
	 * @param file
	 *            a file to be added
	 * @throws IOException
	 */
	public void addFile(final String name, final String fileName,
			final File file) throws IOException {
		try {
			// 去除 除了MINDDUMP的其他信息, 保留注释,备用
			// final StringBuilder stringBuilder = createHeaderBuilder(name);
			// stringBuilder.append("\"; filename=\"").append(fileName)
			// .append("\"\nContent-Type: application/octet-stream\n\n");

			// final String data = stringBuilder.toString();

			// mTotalLength += file.length() + data.length();
			mTotalLength += file.length();
			// mInputChuncks.add(new ByteArrayInputStream(data.getBytes()));
			mInputChuncks.add(new FileInputStream(file));

		} catch (final IOException e) {
			Log.e("TAG", "Can't use input file " + fileName, e);
			throw e;
		}
	}

	public void addStream(String str) {
		byte bs[] = str.getBytes();
		mTotalLength += bs.length;
		mInputChuncks.add(new ByteArrayInputStream(bs));
	}

	public String getUUID() {
		return BOUNDARY_TAG;
	}

	private StringBuilder createHeaderBuilder(final String name) {
		final StringBuilder stringBuilder = new StringBuilder();
		// 去除 除了MINDDUMP的其他信息, 保留注释,备用
		// stringBuilder.append("\n--").append(BOUNDARY_TAG);
		// stringBuilder.append("\nContent-Disposition: form-data; name=\"").append(name);
		stringBuilder.append(name);
		return stringBuilder;
	}

	/**
	 * Finish a body of a post
	 */
	public void finish() {
		Log.w("MultipartHttpEntity", "finish()");

		final String data = "\n--" + BOUNDARY_TAG + "--\n";
		mTotalLength += data.length();
		mInputChuncks.add(new ByteArrayInputStream(data.getBytes()));

		mReady = true;
	}

	/**
	 * (non-Javadoc)
	 * 
	 * @see org.apache.http.HttpEntity#consumeContent()
	 */
	@Override
	public void consumeContent() {
		Log.w("MultipartHttpEntity", "consumeContent()");
		mTotalLength = 0;
		mInputChuncks.clear();

		mReady = false;
	}

	/**
	 * (non-Javadoc)
	 * 
	 * @see org.apache.http.HttpEntity#getContent()
	 */
	@Override
	public InputStream getContent() {
		Log.w("MultipartHttpEntity", "getContent()");
		return new SequenceInputStream(Collections.enumeration(mInputChuncks));
	}

	/**
	 * (non-Javadoc)
	 * 
	 * @see org.apache.http.HttpEntity#getContentEncoding()
	 */
	@Override
	public Header getContentEncoding() {
		Log.w("MultipartHttpEntity", "getContentEncoding()");
		return null;
	}

	/**
	 * (non-Javadoc)
	 * 
	 * @see org.apache.http.HttpEntity#getContentLength()
	 */
	@Override
	public long getContentLength() {
		Log.w("MultipartHttpEntity", "getContentLength() = " + mTotalLength);
		return mTotalLength;
	}

	/**
	 * (non-Javadoc)
	 * 
	 * @see org.apache.http.HttpEntity#getContentType()
	 */
	@Override
	public Header getContentType() {
		Log.w("MultipartHttpEntity", "getContentType()");
		return new BasicHeader("Content-Type", "multipart/form-data; boundary="
				+ BOUNDARY_TAG);
	}

	/**
	 * (non-Javadoc)
	 * 
	 * @see org.apache.http.HttpEntity#isChunked()
	 */
	@Override
	public boolean isChunked() {
		Log.w("MultipartHttpEntity", "isChunked()");
		return false;
	}

	/**
	 * (non-Javadoc)
	 * 
	 * @see org.apache.http.HttpEntity#isRepeatable()
	 */
	@Override
	public boolean isRepeatable() {
		Log.w("MultipartHttpEntity", "isRepeatable()");
		return false;
	}

	/**
	 * (non-Javadoc)
	 * 
	 * @see org.apache.http.HttpEntity#isStreaming()
	 */
	@Override
	public boolean isStreaming() {
		Log.w("MultipartHttpEntity", "isStreaming() = " + mReady);
		return mReady;
	}

	/**
	 * (non-Javadoc)
	 * 
	 * @see org.apache.http.HttpEntity#writeTo(OutputStream)
	 */
	@Override
	public void writeTo(final OutputStream outstream) {
		Log.w("MultipartHttpEntity", "writeTo()");
		long current = 0;
		for (final InputStream inp : mInputChuncks) {
			current += writeFromInputToOutput(inp, outstream, current);
		}
		Log.w("MultipartHttpEntity", "writeEnd()");
	}

	private long writeFromInputToOutput(final InputStream source,
			final OutputStream dest, long current) {
		final byte[] buffer = new byte[BUFFER_SIZE];
		int bytesRead = EOF_MARK;
		int count = 0;
		try {
			while ((bytesRead = source.read(buffer)) != EOF_MARK) {
				Log.w("MultipartHttpEntity", "read = " + bytesRead);

				dest.write(buffer, 0, bytesRead);
				count += bytesRead;

				if (progressCallback != null) {
					progressCallback.onProgress(current + count, mTotalLength);
				}
			}

		} catch (final IOException e) {
			Log.e("TAG", "IOException", e);
		}
		return (long) count;
	}
}
