// Package files implements the chunked uploader and downloader used
// by the account client. The wire protocol requires:
//
// - an initial `requestSendFile` call to obtain a (file_id, dc_id,
// upload_url, access_hash_send) tuple.
// - the actual file is then POSTed to `upload_url` in chunks of
// `chunkSize` bytes, with custom headers per chunk.
// - the server returns `OK/OK` when the upload is complete; if it
// returns `ERROR_TRY_AGAIN`, the upload is restarted.
//
// Downloads work similarly: a per-chunk POST to
// `https://messenger{dc_id}.iranlms.ir/GetFile.ashx` returns the raw
// bytes, optionally in parallel.
package files

// ProgressCallback is invoked after each chunk upload/download with
// the current (total, current) byte counters.
type ProgressCallback func(total, current int64)
