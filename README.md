# AI Gateway — Phase 1 (Raw Proxy)

این نسخه اولین گام از پروژه‌ی AI Gateway است: یک proxy خام و بدون کش/صف بین کلاینت و OpenAI.
هدف این فاز صرفاً این است که مسیر انتها‌به‌انتها (end-to-end) کار کند، قبل از افزودن پیچیدگی‌های بعدی.

## ساختار پروژه

```
ai-gateway/
├── cmd/server/main.go          # نقطه‌ی ورود برنامه، wiring همه‌چیز به هم
├── internal/
│   ├── config/config.go        # بارگذاری تنظیمات از env vars
│   ├── provider/
│   │   ├── provider.go         # interface اصلی LLMProvider + مدل‌های داده
│   │   └── openai.go           # پیاده‌سازی provider برای OpenAI
│   └── handler/chat.go         # HTTP handler برای /v1/chat/completions
└── go.mod
```

## چرا این ساختار؟

- **`internal/provider`** جدا از handler است چون در فاز ۴ باید بتوانیم چند provider
  (OpenAI, Anthropic, ...) را بدون تغییر در handler اضافه کنیم. کافی است یک
  struct جدید بسازیم که interface را برآورده کند.
- **`internal/handler`** فقط مسئول HTTP است (parse request, validate, نوشتن response)؛
  منطق ارتباط با provider را نمی‌داند و مستقیم به آن وابسته نیست، بلکه به interface وابسته است
  (dependency inversion) — همین باعث می‌شود تست‌نویسی handler با یک provider fake ساده باشد.
- **`cmd/server/main.go`** تنها جایی است که همه‌چیز به هم wire می‌شود.

## اجرا

```bash
export OPENAI_API_KEY=sk-...
go run ./cmd/server
```

سپس تست:

```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-4o-mini",
    "messages": [{"role": "user", "content": "سلام"}]
  }'
```

## فاز بعدی

فاز ۲: افزودن Redis caching — قبل از فراخوانی provider، چک می‌کنیم آیا
همین درخواست (بر اساس هش model+messages+temperature) قبلاً پاسخ داده شده.
