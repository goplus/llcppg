#include "lua.h"

int(lua_error)(lua_State *L) {}
void(lua_concat)(lua_State *L, int n) {}
int(lua_next)(lua_State *L, int idx) {}
void(lua_len)(lua_State *L, int idx) {}
long unsigned int(lua_stringtonumber)(lua_State *L, const char *s) {}
void(lua_setallocf)(lua_State *L, lua_Alloc f, void *ud) {}
void(lua_toclose)(lua_State *L, int idx) {}
void(lua_closeslot)(lua_State *L, int idx) {}
