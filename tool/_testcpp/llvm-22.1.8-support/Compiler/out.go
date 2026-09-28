package Compiler

const LLGoPackage = "link: -L$(llvm-config --libdir) -lLLVM; -lLLVM"
const MEMORY_SANITIZER_BUILD = 0
const ADDRESS_SANITIZER_BUILD = 0
const HWADDRESS_SANITIZER_BUILD = 0
const THREAD_SANITIZER_BUILD = 0
const ENABLE_EXCEPTIONS = 1
