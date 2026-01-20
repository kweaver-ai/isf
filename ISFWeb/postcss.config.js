const path = require('path')

module.exports = {
    plugins: [
        ['postcss-import', {
            resolve: function (id, basedir, importOptions) {
                if (id.startsWith('@/')) {
                    /**
                     * 配置别名 @/
                     * 例如 import '@/core/style/base.css' 会被解析为 import '/src/core/style/base.css'
                     */
                    return path.resolve(__dirname, './src', id.substr(2))
                }

                return id
            }
        }],
        'postcss-preset-env',
        'postcss-global-import',
    ]
}