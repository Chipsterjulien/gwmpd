/* A polyfill for browsers that don't support ligatures. */
/* The script tag referring to this file must be placed before the ending body tag. */

/* To provide support for elements dynamically added, this script adds
   method 'icomoonLiga' to the window object. You can pass element references to this method.
*/
(function () {
    'use strict';
    function supportsProperty(p) {
        var prefixes = ['Webkit', 'Moz', 'O', 'ms'],
            i,
            div = document.createElement('div'),
            ret = p in div.style;
        if (!ret) {
            p = p.charAt(0).toUpperCase() + p.substr(1);
            for (i = 0; i < prefixes.length; i += 1) {
                ret = prefixes[i] + p in div.style;
                if (ret) {
                    break;
                }
            }
        }
        return ret;
    }
    var icons;
    if (!supportsProperty('fontFeatureSettings')) {
        icons = {
            'add': '&#xe145;',
            'autorenew': '&#xe863;',
            'create': '&#xe254;',
            'edit': '&#xe254;',
            'mode_edit': '&#xe254;',
            'keyboard_arrow_down': '&#xe313;',
            'keyboard_arrow_left': '&#xe314;',
            'keyboard_arrow_right': '&#xe315;',
            'keyboard_arrow_up': '&#xe316;',
            'loop': '&#xe627;',
            'sync': '&#xe627;',
            'pause': '&#xe034;',
            'play_arrow': '&#xe037;',
            'power_settings_new': '&#xe8ac;',
            'public': '&#xe80b;',
            'queue_music': '&#xe03d;',
            'remove': '&#xe15b;',
            'settings': '&#xe8b8;',
            'skip_next': '&#xe044;',
            'skip_previous': '&#xe045;',
            'stop': '&#xe047;',
            'volume_down': '&#xe04d;',
            'volume_mute': '&#xe04e;',
            'volume_off': '&#xe04f;',
            'volume_up': '&#xe050;',
            'whatshot': '&#xe80e;',
          '0': 0
        };
        delete icons['0'];
        window.icomoonLiga = function (els) {
            var classes,
                el,
                i,
                innerHTML,
                key;
            els = els || document.getElementsByTagName('*');
            if (!els.length) {
                els = [els];
            }
            for (i = 0; ; i += 1) {
                el = els[i];
                if (!el) {
                    break;
                }
                classes = el.className;
                if (/icon-/.test(classes)) {
                    innerHTML = el.innerHTML;
                    if (innerHTML && innerHTML.length > 1) {
                        for (key in icons) {
                            if (icons.hasOwnProperty(key)) {
                                innerHTML = innerHTML.replace(new RegExp(key, 'g'), icons[key]);
                            }
                        }
                        el.innerHTML = innerHTML;
                    }
                }
            }
        };
        window.icomoonLiga();
    }
}());
