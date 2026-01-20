declare namespace Core {
    namespace APIs {
        /**
         * 客户端API调用函数
         * @param P 请求参数
         * @param R 响应内容
         */
        type ThriftAPI<P, R> = (params?: P) => Promise<Readonly<R>>

    }
}