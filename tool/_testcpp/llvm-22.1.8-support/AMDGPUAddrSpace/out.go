package AMDGPUAddrSpace

const LLGoPackage = "link: -L$(llvm-config --libdir) -lLLVM; -lLLVM"
const (
	AMDGPUAS_MAX_AMDGPU_ADDRESS = 9
	// < Address space for flat memory.
	AMDGPUAS_FLAT_ADDRESS = 0
	// < Address space for global memory (RAT0, VTX0).
	AMDGPUAS_GLOBAL_ADDRESS = 1
	// < Address space for region memory. (GDS)
	AMDGPUAS_REGION_ADDRESS = 2
	// < Address space for local memory.
	AMDGPUAS_LOCAL_ADDRESS = 3
	// < Address space for constant memory (VTX2).
	AMDGPUAS_CONSTANT_ADDRESS = 4
	// < Address space for private memory.
	AMDGPUAS_PRIVATE_ADDRESS = 5
	// < Address space for 32-bit constant memory.
	AMDGPUAS_CONSTANT_ADDRESS_32BIT = 6
	// < Address space for 160-bit buffer fat pointers.
	// < Not used in backend.
	AMDGPUAS_BUFFER_FAT_POINTER = 7
	// < Address space for 128-bit buffer resources.
	AMDGPUAS_BUFFER_RESOURCE = 8
	// < Address space for 192-bit fat buffer
	// < pointers with an additional index.
	AMDGPUAS_BUFFER_STRIDED_POINTER = 9
	// < Address space for GS NGG Streamout registers.
	AMDGPUAS_STREAMOUT_REGISTER = 128
	// Address space for direct addressable parameter memory (CONST0).
	AMDGPUAS_PARAM_D_ADDRESS = 6
	// Address space for indirect addressable parameter memory (VTX1).
	AMDGPUAS_PARAM_I_ADDRESS = 7
	// Address space for indirect addressable parameter memory (VTX1).
	AMDGPUAS_CONSTANT_BUFFER_0 = 8
	// Address space for indirect addressable parameter memory (VTX1).
	AMDGPUAS_CONSTANT_BUFFER_1 = 9
	// Address space for indirect addressable parameter memory (VTX1).
	AMDGPUAS_CONSTANT_BUFFER_2 = 10
	// Address space for indirect addressable parameter memory (VTX1).
	AMDGPUAS_CONSTANT_BUFFER_3 = 11
	// Address space for indirect addressable parameter memory (VTX1).
	AMDGPUAS_CONSTANT_BUFFER_4 = 12
	// Address space for indirect addressable parameter memory (VTX1).
	AMDGPUAS_CONSTANT_BUFFER_5 = 13
	// Address space for indirect addressable parameter memory (VTX1).
	AMDGPUAS_CONSTANT_BUFFER_6 = 14
	// Address space for indirect addressable parameter memory (VTX1).
	AMDGPUAS_CONSTANT_BUFFER_7 = 15
	// Address space for indirect addressable parameter memory (VTX1).
	AMDGPUAS_CONSTANT_BUFFER_8 = 16
	// Address space for indirect addressable parameter memory (VTX1).
	AMDGPUAS_CONSTANT_BUFFER_9 = 17
	// Address space for indirect addressable parameter memory (VTX1).
	AMDGPUAS_CONSTANT_BUFFER_10 = 18
	// Address space for indirect addressable parameter memory (VTX1).
	AMDGPUAS_CONSTANT_BUFFER_11 = 19
	// Address space for indirect addressable parameter memory (VTX1).
	AMDGPUAS_CONSTANT_BUFFER_12 = 20
	// Address space for indirect addressable parameter memory (VTX1).
	AMDGPUAS_CONSTANT_BUFFER_13 = 21
	// Address space for indirect addressable parameter memory (VTX1).
	AMDGPUAS_CONSTANT_BUFFER_14 = 22
	// Address space for indirect addressable parameter memory (VTX1).
	AMDGPUAS_CONSTANT_BUFFER_15 = 23
	// Address space for indirect addressable parameter memory (VTX1).
	AMDGPUAS_UNKNOWN_ADDRESS_SPACE = -1
)
const (
	AMDGPU_DWARFAS_GLOBAL       = 0
	AMDGPU_DWARFAS_GENERIC      = 1
	AMDGPU_DWARFAS_REGION       = 2
	AMDGPU_DWARFAS_LOCAL        = 3
	AMDGPU_DWARFAS_PRIVATE_LANE = 5
	AMDGPU_DWARFAS_PRIVATE_WAVE = 6
	AMDGPU_DWARFAS_DEFAULT      = 0
)
