# Import and Convert All Conan C Libraries

Use the [rxsync](https://github.com/goplus/rxsync) tool to convert 1,800 libraries from the [Conan Center Index](https://github.com/conan-io/conan-center-index/tree/master/recipes).

## Steps

1. Use `rxsync` to sync Conan recipes locally
2. Customize names and generate a name mapping table
3. Parse dependency information to obtain the conversion list
4. Perform conversion according to the order in the conversion list
5. Automatically populate dependencies and remove unnecessary `include`s based on Conan dependency info
6. After converting each library, submit a pull request. Use GitHub Actions to perform a basic validation via `llgo build`
7. After merging, tag each `{{clib}}/{{cver}}` folder as `{{clib}}/{{cver}}/v0.1.0` to meet nested module requirements

## Issues

### 1. All `cflags` Included by Default

The automatically generated `llcppg.cfg` includes all header files from dependencies listed under `cflags`. For example, `libxml2` includes header files from `zlib`, and `libxslt` includes headers from both `libxml2` and `zlib` as part of its include paths. This behavior is unintended — for example, in certain cases, headers from `libxml2` and `zlib` need to be removed.

**Preliminary solution:** Automatically remove them during dependency population.

### 2. Local `deps`

How should dependencies be populated for `deps` in `llcppg.cfg`? All batch processing is done locally, while `deps` is expected to reference the GitHub path under `github.com/llpkg`. There's uncertainty about whether locally converted packages can be referenced as dependencies during batch processing.

**Needs verification.**

### 3. Package Order

Sorted by popularity of repository and dependency.

Rule:
1. Score(Stars/Forks/Watchers/Contributors)
2. The number of dependency(Less dependency is higher)

Tools: [Rank Tool](https://github.com/MeteorsLiu/reporank)

### 4. Naming Conflicts in Custom Names

Remove `lib` prefix.

However, this may cause naming conflicts and needs to be verified.

### 5. Demo Validation

Batch conversion cannot perform demo validation.

**Preliminary solution:** Use basic `llgo build` validation.

### 6. Trim Prefix

The expected behavior is to configure `trimPrefixes` in `llcppg.cfg` to remove namespaces before committing to the repo.
However, the automated process may not be able to identify and fill this in automatically.

**No solution yet.**

### 7. Tagging for All Version Branches

How should tags be created?

**Preliminary solution:**
```
cjson/3.4.5/v0.1.0
cjson/3.3.5/v0.1.0
cjson/3.2.5/v0.1.0
```

But it's unclear whether Go can recognize this format — needs verification.

### 8. GitHub Action Verification

Is it necessary to add `llcppg`-based validation?
