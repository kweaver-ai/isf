/* eslint-disable */
const withLess = require("@zeit/next-less");
const path = require("path");
const isDEV = process.env.NODE_ENV === "development";
const prdConfig = isDEV
  ? {}
  : {
      webpack(config, { buildId }) {
        config.module.rules[0].include.push(
          /templite[\\/]dist/,
          /rosetta[\\/]dist/
        );
        const { exclude } = config.module.rules[0];
        config.module.rules[0].exclude = (excludePath) => {
          if (
            [/templite[\\/]dist/, /rosetta[\\/]dist/].some((reg) =>
              reg.test(excludePath)
            )
          ) {
            return false;
          }
          return exclude(excludePath);
        };
        return config;
      },
    };

// SVG 加载器配置
const svgrLoaderConfig = {
  test: /\.svg$/,
  use: [
    {
      loader: "@svgr/webpack",
      options: {
        typescript: false,
        svgoConfig: {
          plugins: {
            removeViewBox: false,
          },
        },
      },
    },
  ],
};

// 配置别名
const aliasConfig = {
  "@icons": path.resolve(__dirname, "src/icons"),
};

module.exports = withLess({
  transpileModules: ["templite", "rosetta"],
  lessLoaderOptions: {
    lessOptions: {
      javascriptEnabled: true,
    },
  },
  typescript: {
    ignoreDevErrors: true,
  },
  basePath: "/oauth2-ui",
  experimental: {
    basePath: "/oauth2-ui",
  },
  assetPrefix: "./",
  ...prdConfig,
  webpack(config, options) {
    config.module.rules.push(svgrLoaderConfig);

    // 添加处理图片文件的规则
    config.module.rules.push({
      test: /\.(png|jpe?g|gif)$/i,
      use: [
        {
          loader: "file-loader",
          options: {
            publicPath: "/oauth2/_next/static/images/",
            outputPath: `/static/images/`,
            name: "[name].[hash].[ext]",
          },
        },
      ],
    });

    config.resolve.alias = {
      ...config.resolve.alias,
      ...aliasConfig,
    };
    return config;
  },
});
