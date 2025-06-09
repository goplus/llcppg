package bzip3

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

const OK = 0

type State struct {
	Unused [8]uint8
}

/**
 * @brief Get bzip3 version.
 */
//go:linkname Version C.bz3_version
func Version() *c.Char

/**
 * @brief Get the last error number associated with a given state.
 */
// llgo:link (*State).LastError C.bz3_last_error
func (recv_ *State) LastError() c.Int8T {
	return 0
}

/**
 * @brief Return a user-readable message explaining the cause of the last error.
 */
// llgo:link (*State).Strerror C.bz3_strerror
func (recv_ *State) Strerror() *c.Char {
	return nil
}

/**
 * @brief Construct a new block encoder state, which will encode blocks as big as the given block size.
 * The decoder will be able to decode blocks at most as big as the given block size.
 * Returns NULL in case allocation fails or the block size is not between 65K and 511M
 */
//go:linkname New C.bz3_new
func New(block_size c.Int32T) *State

/**
 * @brief Free the memory occupied by a block encoder state.
 */
// llgo:link (*State).Free C.bz3_free
func (recv_ *State) Free() {
}

/**
 * @brief Return the recommended size of the output buffer for the compression functions.
 */
//go:linkname Bound C.bz3_bound
func Bound(input_size c.SizeT) c.SizeT

/**
 * @brief Compress a frame. This function does not support parallelism
 * by itself, consider using the low level `bz3_encode_blocks()` function instead.
 * Using the low level API might provide better performance.
 * Returns a bzip3 error code; BZ3_OK when the operation is successful.
 * Make sure to set out_size to the size of the output buffer before the operation;
 * out_size must be at least equal to `bz3_bound(in_size)'.
 */
//go:linkname Compress C.bz3_compress
func Compress(block_size c.Uint32T, in *c.Uint8T, out *c.Uint8T, in_size c.SizeT, out_size *c.SizeT) c.Int

/**
 * @brief Decompress a frame. This function does not support parallelism
 * by itself, consider using the low level `bz3_decode_blocks()` function instead.
 * Using the low level API might provide better performance.
 * Returns a bzip3 error code; BZ3_OK when the operation is successful.
 * Make sure to set out_size to the size of the output buffer before the operation.
 */
//go:linkname Decompress C.bz3_decompress
func Decompress(in *c.Uint8T, out *c.Uint8T, in_size c.SizeT, out_size *c.SizeT) c.Int

/**
 * @brief Calculate the minimal memory required for compression with the given block size.
 * This includes all internal buffers and state structures. This calculates the amount of bytes
 * that will be allocated by a call to `bz3_new()`.
 *
 * @details Memory allocation and usage patterns:
 *
 * bz3_new():
 *    - Allocates all memory upfront:
 *      - Core state structure (sizeof(struct bz3_state))
 *      - Swap buffer (bz3_bound(block_size) bytes)
 *      - SAIS array (BWT_BOUND(block_size) * sizeof(int32_t) bytes)
 *      - LZP lookup table ((1 << LZP_DICTIONARY) * sizeof(int32_t) bytes)
 *      - Compression state (sizeof(state))
 *    - All memory remains allocated until bz3_free()
 *
 * Additional memory may be used depending on API used from here.
 *
 * # Low Level APIs
 *
 * 1. bz3_encode_block() / bz3_decode_block():
 *    - Uses pre-allocated memory from bz3_new()
 *    - No additional memory allocation except for libsais (usually ~16KiB)
 *    - Peak memory usage of physical RAM varies with compression stages:
 *      - LZP: Uses LZP lookup table + swap buffer
 *      - BWT: Uses SAIS array + swap buffer
 *      - Entropy coding: Uses compression state (cm_state) + swap buffer
 *
 * Using the higher level API, `bz3_compress`, expect an additional allocation
 * of `bz3_bound(block_size)`.
 *
 * In the parallel version `bz3_encode_blocks`, each thread gets its own state,
 * so memory usage is `n_threads * bz3_compress_memory_needed()`.
 *
 * # High Level APIs
 *
 * 1. bz3_compress():
 *    - Allocates additional temporary compression buffer (bz3_bound(block_size) bytes)
 *      in addition to the memory amount returned by this method call and libsais.
 *    - Everything is freed after compression completes
 *
 * 2. bz3_decompress():
 *    - Allocates additional temporary compression buffer (bz3_bound(block_size) bytes)
 *      in addition to the memory amount returned by this method call and libsais.
 *    - Everything is freed after compression completes
 *
 * Memory remains constant during operation, with except of some small allocations from libsais during
 * BWT stage. That is not accounted by this function, though it usually amounts to ~16KiB, negligible.
 * The worst case of BWT is 2*block_size technically speaking.
 *
 * No dynamic (re)allocation occurs outside of that.
 *
 * @param block_size The block size to be used for compression
 * @return The total number of bytes required for compression, or 0 if block_size is invalid
 */
//go:linkname MinMemoryNeeded C.bz3_min_memory_needed
func MinMemoryNeeded(block_size c.Int32T) c.SizeT

/**
 * @brief Encode a single block. Returns the amount of bytes written to `buffer'.
 * `buffer' must be able to hold at least `bz3_bound(size)' bytes. The size must not
 * exceed the block size associated with the state.
 */
// llgo:link (*State).EncodeBlock C.bz3_encode_block
func (recv_ *State) EncodeBlock(buffer *c.Uint8T, size c.Int32T) c.Int32T {
	return 0
}

/**
 * @brief Decode a single block.
 *
 * `buffer' must be able to hold at least `bz3_bound(orig_size)' bytes
 * in order to ensure decompression will succeed for all possible bzip3 blocks.
 *
 * In most (but not all) cases, `orig_size` should usually be sufficient.
 * If it is not sufficient, you must allocate a buffer of size `bz3_bound(orig_size)` temporarily.
 *
 * If `buffer_size` is too small, `BZ3_ERR_DATA_SIZE_TOO_SMALL` will be returned.
 * The size must not exceed the block size associated with the state.
 *
 * @param buffer_size The size of the buffer at 'buffer'
 * @param compressed_size The size of the compressed data in 'buffer'
 * @param orig_size The original size of the data before compression.
 */
// llgo:link (*State).DecodeBlock C.bz3_decode_block
func (recv_ *State) DecodeBlock(buffer *c.Uint8T, buffer_size c.SizeT, compressed_size c.Int32T, orig_size c.Int32T) c.Int32T {
	return 0
}

/**
 * @brief Check if using original file size as buffer size is sufficient for decompressing
 * a block at `block` pointer.
 *
 * @param block Pointer to the compressed block data
 * @param block_size Size of the block buffer in bytes (must be at least 13 bytes for header)
 * @param orig_size Size of the original uncompressed data
 * @return 1 if original size is sufficient, 0 if insufficient, -1 on header error (insufficient buffer size)
 *
 * @remarks
 *
 *      This function is useful for external APIs using the low level block encoding API,
 *      `bz3_encode_block`. You would normally call this directly after `bz3_encode_block`
 *      on the block that has been output.
 *
 *      The purpose of this function is to prevent encoding blocks that would require an additional
 *      malloc at decompress time.
 *      The goal is to prevent erroring with `BZ3_ERR_DATA_SIZE_TOO_SMALL`, thus
 *      in turn
 */
//go:linkname OrigSizeSufficientForDecode C.bz3_orig_size_sufficient_for_decode
func OrigSizeSufficientForDecode(block *c.Uint8T, block_size c.SizeT, orig_size c.Int32T) c.Int
