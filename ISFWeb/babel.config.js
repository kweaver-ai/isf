module.exports = {
    "presets": [
        "@babel/preset-react",
        "@babel/preset-typescript",
        [
            "@babel/preset-env",
            {
                "useBuiltIns": "entry",
                "corejs": "3",
                "modules": false
            }
        ]
    ],
    "plugins": [
        [
            "import",
            {
                "libraryName": "antd",
                "libraryDirectory": "es",
                "style": false,
            },
            "antd"
        ]
    ]
}