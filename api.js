!(function() {
    const API_DATA = {
        prefix: "$$PREFIX$$",
        headers: [],
    };

    function generateApi(name) {
        return async (dat={}) => {
            const response = await fetch(API_DATA.prefix + name, {
                method: "POST",
                headers: API_DATA.headers,
                body: JSON.stringify(dat),
            });
            if (response.status === 200) {
                return await response.json();
            } else {
                throw new Error("Status: " + response.status + "  Error: " + await response.text())
            }
        }
    }

    window.$$API_NAME$$ = new Proxy(API_DATA, {
        get(target, name) {
            const value = target[name]
            if (value !== undefined) return value;
            return generateApi(name);
        },
        set(target, param, value) {
            if (value === undefined) return;
            if (target[param] !== undefined) {
                target[param] = value;
            }
        }
    })
})()

