export class I18N {
    static i(key, param1 ) {
        return $.i18n(key, param1);
    }

    static i18n(key, param1) {
        return $.i18n(key, param1);
    }

    static promise() {
        let locale = navigator.language || navigator.userLanguage;
        $.i18n({
            locale: locale
        });
        console.log('I18N.promise', locale);
        return new Promise(function (resolve) {
            $.i18n().load(`/static/i18n/${locale}.json`, locale)
                .done(function() {
                    resolve()
                });
        })
    }
}
