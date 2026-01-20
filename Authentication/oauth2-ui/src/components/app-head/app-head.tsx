import React, { FunctionComponent, useMemo } from "react";
import { useObserver } from "mobx-react-lite";
import Head from "next/head";
import { generate } from "@ant-design/colors";

interface IAppHeadProps {
  oemconfig: {
    [key: string]: any;
  };
}
export const AppHead: FunctionComponent<IAppHeadProps> = ({ oemconfig }) => {
  const colorPalette: string[] = useMemo(
    () => generate(oemconfig.theme || "#6775CD"),
    [oemconfig]
  );
  return useObserver(() => {
    return (
      <Head>
        <meta
          name="viewport"
          content="width=device-width, initial-scale=1, maximum-scale=1, minimum-scale=1,user-scalable=no"
        />
        <meta http-equiv="X-UA-Compatible" content="IE=11" />
        <title>{oemconfig.product}</title>
        <link
          rel="shortcut icon"
          type="image/x-icon"
          href={`data:image/png;base64,${oemconfig.favicon}`}
        ></link>
        {/* <link rel="stylesheet" type="text/css" href="/static/SourceHanSansCNNormalFont/fontblock.css" /> */}
        <style type="text/css">
          {`
                        /* 兼容Chromium 87的选择器(87不支持:not) - 使用级联覆盖方式 */
                        body * {
                            color: ${
                              oemconfig.isTransparentBoxStyle
                                ? "#ffffff !important"
                                : "#A7AEB8"
                            };
                        }
                        
                        /* 确保链接和按钮内文本保持原始颜色 */
                        body a,
                        body button span {
                            color: inherit !important;
                        }

                        body {
                            background: ${
                              oemconfig.isTransparentBoxStyle
                                ? "transparent"
                                : "#ffffff"
                            } !important;
                        }

                        input {
                            color: ${
                              oemconfig.isTransparentBoxStyle
                                ? "#ffffff !important"
                                : "#000000"
                            };

                            border-bottom-color: ${
                              oemconfig.isTransparentBoxStyle
                                ? "rgba(255, 255, 255, 0.65) !important"
                                : "#e2e7F1"
                            };
                        }

                        input::placeholder {
                          color: ${
                            oemconfig.isTransparentBoxStyle
                              ? "rgba(255, 255, 255, 0.65)"
                              : "#c0c0c0"
                          } !important;
                        }

                        input:-ms-input-placeholder {
                          color: ${
                            oemconfig.isTransparentBoxStyle
                              ? "rgba(255, 255, 255, 0.65)"
                              : "#c0c0c0"
                          } !important;
                        }

                        input::-ms-input-placeholder {
                           color: ${
                             oemconfig.isTransparentBoxStyle
                               ? "rgba(255, 255, 255, 0.65)"
                               : "#c0c0c0"
                           } !important;
                        }

                        .input-item {
                            border-bottom-color: ${
                              oemconfig.isTransparentBoxStyle
                                ? "rgba(255, 255, 255, 0.65) !important"
                                : "#e2e7F1"
                            };
                        }
                            
                        .as-components-oem-background-color {
                            background: ${colorPalette[5]} !important;
                            color: #ffffff !important;
                        }
                        .as-components-oem-background-color:hover {
                            background: ${colorPalette[4]} !important;
                        }
                        .as-components-oem-background-color:active {
                            background: ${colorPalette[6]} !important;
                        }
                        .as-components-oem-background-color[disabled], .as-components-oem-background-color.disabled {
                            background: ${colorPalette[3]} !important;
                            color: #ffffff !important;
                        }
                    `}
        </style>
      </Head>
    );
  });
};
