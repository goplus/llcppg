#define INI_API __attribute__((visibility("default")))
class INIReader
{
public:
    __attribute__((visibility("default"))) explicit INIReader(const char *filename);
    INI_API explicit INIReader(const char *buffer, long buffer_size);
    ~INIReader();
    INI_API int ParseError() const;
    INI_API const char *Get(const char *section, const char *name,
                            const char *default_value) const;

private:
    static const char *MakeKey(const char *section, const char *name);
};