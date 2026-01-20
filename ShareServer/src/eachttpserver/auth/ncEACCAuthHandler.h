#ifndef __NC_EAC_CAUTH_HANDLER_H__
#define __NC_EAC_CAUTH_HANDLER_H__

class ncEACCAuthHandler
{
public:
    ncEACCAuthHandler ();
    ~ncEACCAuthHandler (void);

    void doCARequestHandler (brpc::Controller* cntl);

protected:
    /***
     * 获取管理文档信息
     */
    void Get (brpc::Controller* cntl, const String& userId);

private:
    typedef void (ncEACCAuthHandler::*ncMethodFunc) (brpc::Controller*, const String&);
    map<String, ncMethodFunc>            _methodFuncs;
};

#endif  // __NC_EAC_CAUTH_HANDLER_H__
