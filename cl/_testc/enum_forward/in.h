enum clnt_stat;
enum clnt_stat;

enum clnt_stat *rpc_status(void);

enum clnt_stat {
    RPC_SUCCESS = 0,
    RPC_FAILED = 1,
};
