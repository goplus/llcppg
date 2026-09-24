enum Color {
	Red,
	Green,
	Blue
};

enum State {
	StateStopped = 0,
	StateRunning = 10,
	StatePaused
};

enum {
	FlagA = 1 << 0,
	FlagB = 1 << 1,
	FlagC = 1 << 2
};

/**
 * Error codes for Compilation Database
 */
typedef enum  {
  /*
   * No error occurred
   */
  CXCompilationDatabase_NoError = 0,

  /*
   * Database can not be loaded
   */
  CXCompilationDatabase_CanNotLoadDatabase = 1

} CXCompilationDatabase_Error;

void f(int, CXCompilationDatabase_Error);
