#include <stdbool.h>

#if defined(__cplusplus)

#include <functional>

class deferred_expr
{
public:
    deferred_expr(std::function<void()> f_ptr) : _f_ptr(f_ptr)
    {
    }
    ~deferred_expr()
    {
        if (_f_ptr == nullptr)
            return;
        _f_ptr();
    }
    void reset()
    {
        _f_ptr = nullptr;
    }

private:
    std::function<void()> _f_ptr;
};

#define _D_MERGE_(a, b) a##b
#define _D_MERGE(a, b) _D_MERGE_(a, b)

#define defer(lmd) auto _D_MERGE(__defer_at_, __LINE__) = deferred_expr(lmd)

extern "C" {
#endif

bool open_directory(const char *title, char *out_path);
bool open_file(const char *title, char *out_path);
bool save_file(const char *title, char *out_path);

#if defined(__cplusplus)
}
#endif