## Описание
Данный discord bot реализует следующие команды:
- !help - выводит описание всех доступных команд;
- !weather - выводит информацию о погоде по переданному населенному пункту [OpenWeatherMap](https://openweathermap.org/);
- !translate - переводит на выбранный язык переданный текст через [Google translate](https://translate.google.com/);
- !translate_lang - выводит список доступных языков для перевода.

Бот со всеми компонентами можно развернуть в Docker через docker-compose up.  
Для корректной работы с API необходимо прописать в файле bot.env ключи
   
    DISCORD_TOKEN="Discord bot key"
    OWM_KEY="OpenWeatherMap API key"
    GOOGLE_KEY="Google API key"

Ссылка для [добавления бота](https://discord.com/api/oauth2/authorize?client_id=1202914737793278022&permissions=8&scope=bot)  
Примеры:
- !weather
---
    !weather Almaty
    Погода в Алматы
    Условия
    туман
    Температура
    -8.05°C
    Влажность
    73%
    Ветер
    3.00 м/c
---
    !weather Москва
    Погода в Москва
    Условия
    пасмурно
    Температура
    -2.19°C
    Влажность
    96%
    Ветер
    3.14 м/c

- !translate
---
    !translate ru Hello world!
    Привет, мир!
---
    !translate Russian Hello world!
    Привет, мир!
