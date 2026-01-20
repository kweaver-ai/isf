const path = require('path')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const CSSMinimizerWebpackPlugin = require('css-minimizer-webpack-plugin')
const CopyWebpackPlugin = require("copy-webpack-plugin")
const { name: packageName } = require('./package')
const fs = require('fs')

const pages = fs.readdirSync('./src/pages') // such as ['page1.tsx', 'page2.tsx']
const { entry, htmls } =  pages.reduce((prev, page) => { // 多页面打包
    const [filename] = page.split('.')
    const { entry: prevEntry, htmls: prevHtmls } = prev
    const entry = {
        ...prevEntry,
        [filename]: `./src/pages/${page}`,
    }
    const htmls = [
        ...prevHtmls,
        new HtmlWebpackPlugin({
            template: './src/public/index.html',
            filename: `${filename}.html`,
            chunks: [filename],
        })
    ]

    return { entry, htmls }
}, { entry: {}, htmls: [] })

const config = {
    entry,
    output: {
        publicPath: './',
        path: path.resolve(__dirname, './dist'),
        filename: 'scripts/[name].[contenthash].js',
        chunkFilename: 'scripts/[name].chunk.[chunkhash].js',
        library: `${packageName}-[name]`, // 子应用必须
        libraryTarget: 'umd', // 子应用必须
    },
    module: {
        rules: [
            {
                test: /\.jsx?$|\.tsx?$/,
                use: 'babel-loader',
                sideEffects: false, // 便于进行tree-shaking
                exclude: [
                    /antd/,
                ],
            },
            {
                test: /\.css$/,
                use: [
                    {
                        loader: MiniCssExtractPlugin.loader, // 将css文件打包成独立的文件
                    },
                    {
                        loader: 'css-loader',
                        options: {
                            modules: { // 启用css模块
                                localIdentName: '[path][name]---[local]', // 打包的class名称
                            },
                            importLoaders: 1,
                        },
                    },
                    {
                        loader: 'postcss-loader',
                        options: {
                            postcssOptions: {
                                config: path.resolve(__dirname, './postcss.config.js'),
                            },
                        },
                    }
                ]
            },
            {
                test: /\.less$/,
                use: [
                    {
                        loader: 'style-loader',
                    },
                    {
                        loader: 'css-loader',
                    },
                    {
                        loader: 'less-loader',
                        options: {
                            lessOptions: {
                                javascriptEnabled: true,
                            },
                        },
                    }
                ],
            },
            {
                test: /\.png$|\.gif$/,
                use: {
                    loader: 'url-loader',
                    options: {
                        name: 'assets/images/[name]_[hash].[ext]',
                        outputPath: '',
                        limit: 1024 * 1024 * 1024,
                        esModule: false,
                    }
                },
                type: 'javascript/auto',
            },
            {
                test: /\.(eot|otf|webp|svg|ttf|woff|woff2)(\?.*)?$/,
                use: {
                    loader: 'url-loader',
                    options: {
                        name: 'assets/fonts/[name]_[hash].[ext]',
                        outputPath: '',
                        limit: 1024 * 1024 * 1024,
                        esModule: false,
                    }
                },
                type: 'javascript/auto',
                exclude: path.resolve(__dirname, './src/icons'), 
            },
            {
                test: /\.svg$/,
                use: [
                    {
                        loader: "@svgr/webpack",
                        options: {
                            typescript: false,
                            svgoConfig: {
                                plugins: [{
                                    name: 'removeViewBox', 
                                    active: false 
                                }],
                            },
                        },
                    },
                ],
                include: path.resolve(__dirname, './src/icons'), 
            }
        ],
    },
    plugins: [
        ...htmls,
        new CopyWebpackPlugin({
            patterns: [
                {
                    from: path.resolve(__dirname, './src/icons'),
                    to: 'icons'
                },
            ]
        }),
        new MiniCssExtractPlugin({
            filename: 'styles/[name].[contenthash].css',
            ignoreOrder: true,
        }),
        new CopyWebpackPlugin({
            patterns: [
                {
                    from: "node_modules/@dip/components/dist/dip-components.min.css",
                    to: "styles/dip-components.min.css",
                }
            ]
        })
    ],
    resolve: {
        extensions: ['.ts', '.tsx', '.js', '.jsx', '.css'],
        alias: {
            '@': path.resolve(__dirname, './src'),
        },
    },
    optimization: {
        splitChunks: {
            cacheGroups: {
                vendor: { // 提取node_modules到vendor.js
                    test: /node_modules/,
                    name: 'vendor',
                    chunks: 'all',
                    priority: -10,
                },
                common: { // 提取多次导入的模块到common.js（貌似在动态加载的模块才起作用）
                    name: 'common',
                    chunks: 'all', // async（默认） | initial | all。async - 只有在异步加载模块的时候，才会进行分包处理该模块。initial - 同步加载，进行分包处理；all - 同步异步都会进行分包处理
                    minChunks: 2,
                    minSize: 10,
                    priority: -20,
                    reuseExistingChunk: true,
                },
            },
        },
    }
}

module.exports = (env, { mode }) => {
    if (mode === 'development') {
        config.devtool = 'source-map'
        config.devServer = {
            compress: true,
            host: '0.0.0.0',
            port: 1001,
            hot: false,
            headers: {
                'Access-Control-Allow-Origin': '*',
            },
        }
        config.output.publicPath = ''
    } else {
        config.optimization.minimizer =  [
            `...`, // 保留默认的js压缩
            new CSSMinimizerWebpackPlugin(), // 压缩css代码
        ]
    }

    return config
}