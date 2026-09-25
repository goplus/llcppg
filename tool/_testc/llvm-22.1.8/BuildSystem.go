package clang

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

type VirtualFileOverlayImpl struct {
}

// Object encapsulating information about overlaying virtual
// file/directories over the real file system.
type VirtualFileOverlay = *VirtualFileOverlayImpl
type ModuleMapDescriptorImpl struct {
}

// Object encapsulating information about a module.modulemap file.
type ModuleMapDescriptor = *ModuleMapDescriptorImpl

// Return the timestamp for use with Clang's
// \c -fbuild-session-timestamp= option.
//
//go:linkname GetBuildSessionTimestamp C.clang_getBuildSessionTimestamp
func GetBuildSessionTimestamp() c.UlongLong

// Create a \c CXVirtualFileOverlay object.
// Must be disposed with \c clang_VirtualFileOverlay_dispose().
//
// \param options is reserved, always pass 0.
//
//go:linkname VirtualFileOverlayCreate C.clang_VirtualFileOverlay_create
func VirtualFileOverlayCreate(options c.Uint) VirtualFileOverlay

// Map an absolute virtual file path to an absolute real one.
// The virtual path must be canonicalized (not contain "."/"..").
// \returns 0 for success, non-zero to indicate an error.
//
// llgo:link (*VirtualFileOverlayImpl).VirtualFileOverlayAddFileMapping C.clang_VirtualFileOverlay_addFileMapping
func (_llcppg_param1 *VirtualFileOverlayImpl) VirtualFileOverlayAddFileMapping(virtualPath *c.Char, realPath *c.Char) ErrorCode {
	return 0
}

// Set the case sensitivity for the \c CXVirtualFileOverlay object.
// The \c CXVirtualFileOverlay object is case-sensitive by default, this
// option can be used to override the default.
// \returns 0 for success, non-zero to indicate an error.
//
// llgo:link (*VirtualFileOverlayImpl).VirtualFileOverlaySetCaseSensitivity C.clang_VirtualFileOverlay_setCaseSensitivity
func (_llcppg_param1 *VirtualFileOverlayImpl) VirtualFileOverlaySetCaseSensitivity(caseSensitive c.Int) ErrorCode {
	return 0
}

// Write out the \c CXVirtualFileOverlay object to a char buffer.
//
// \param options is reserved, always pass 0.
// \param out_buffer_ptr pointer to receive the buffer pointer, which should be
// disposed using \c clang_free().
// \param out_buffer_size pointer to receive the buffer size.
// \returns 0 for success, non-zero to indicate an error.
//
// llgo:link (*VirtualFileOverlayImpl).VirtualFileOverlayWriteToBuffer C.clang_VirtualFileOverlay_writeToBuffer
func (_llcppg_param1 *VirtualFileOverlayImpl) VirtualFileOverlayWriteToBuffer(options c.Uint, out_buffer_ptr **c.Char, out_buffer_size *c.Uint) ErrorCode {
	return 0
}

// free memory allocated by libclang, such as the buffer returned by
// \c CXVirtualFileOverlay() or \c clang_ModuleMapDescriptor_writeToBuffer().
//
// \param buffer memory pointer to free.
//
//go:linkname Free C.clang_free
func Free(buffer unsafe.Pointer)

// Dispose a \c CXVirtualFileOverlay object.
//
// llgo:link (*VirtualFileOverlayImpl).VirtualFileOverlayDispose C.clang_VirtualFileOverlay_dispose
func (_llcppg_param1 *VirtualFileOverlayImpl) VirtualFileOverlayDispose() {
}

// Create a \c CXModuleMapDescriptor object.
// Must be disposed with \c clang_ModuleMapDescriptor_dispose().
//
// \param options is reserved, always pass 0.
//
//go:linkname ModuleMapDescriptorCreate C.clang_ModuleMapDescriptor_create
func ModuleMapDescriptorCreate(options c.Uint) ModuleMapDescriptor

// Sets the framework module name that the module.modulemap describes.
// \returns 0 for success, non-zero to indicate an error.
//
// llgo:link (*ModuleMapDescriptorImpl).ModuleMapDescriptorSetFrameworkModuleName C.clang_ModuleMapDescriptor_setFrameworkModuleName
func (_llcppg_param1 *ModuleMapDescriptorImpl) ModuleMapDescriptorSetFrameworkModuleName(name *c.Char) ErrorCode {
	return 0
}

// Sets the umbrella header name that the module.modulemap describes.
// \returns 0 for success, non-zero to indicate an error.
//
// llgo:link (*ModuleMapDescriptorImpl).ModuleMapDescriptorSetUmbrellaHeader C.clang_ModuleMapDescriptor_setUmbrellaHeader
func (_llcppg_param1 *ModuleMapDescriptorImpl) ModuleMapDescriptorSetUmbrellaHeader(name *c.Char) ErrorCode {
	return 0
}

// Write out the \c CXModuleMapDescriptor object to a char buffer.
//
// \param options is reserved, always pass 0.
// \param out_buffer_ptr pointer to receive the buffer pointer, which should be
// disposed using \c clang_free().
// \param out_buffer_size pointer to receive the buffer size.
// \returns 0 for success, non-zero to indicate an error.
//
// llgo:link (*ModuleMapDescriptorImpl).ModuleMapDescriptorWriteToBuffer C.clang_ModuleMapDescriptor_writeToBuffer
func (_llcppg_param1 *ModuleMapDescriptorImpl) ModuleMapDescriptorWriteToBuffer(options c.Uint, out_buffer_ptr **c.Char, out_buffer_size *c.Uint) ErrorCode {
	return 0
}

// Dispose a \c CXModuleMapDescriptor object.
//
// llgo:link (*ModuleMapDescriptorImpl).ModuleMapDescriptorDispose C.clang_ModuleMapDescriptor_dispose
func (_llcppg_param1 *ModuleMapDescriptorImpl) ModuleMapDescriptorDispose() {
}
