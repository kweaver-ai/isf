#ifndef __NC_ACS_DATA_STORE_H__
#define __NC_ACS_DATA_STORE_H__

#include <ossclient/public/ncIOSSClient.h>
#include <public/ncIACSDataStore.h>
#include <dataapi/dataapi.h>
#include <boost/thread/tss.hpp>
#include <drivenadapter/public/ossGatewayInterface.h>

class ncDataCopier :public ncIDataCopier
{
public:
    ncDataCopier (string& data)
        : _data (data)
        , _length (0)
    {}

    NS_DECL_ISUPPORTS
    NS_DECL_NCIDATACOPIER

private:
    string& _data;
    int _length;
};

/* Header file */
class ncACSDataStore : public ncIACSDataStore
{
    AB_DECLARE_THREADSAFE_SINGLETON (ncACSDataStore)

public:
  NS_DECL_ISUPPORTS
  NS_DECL_NCIACSDATASTORE

  ncACSDataStore();
  void getAccountObjId (const String& id, String& accountId, String& objId);

private:
  ~ncACSDataStore();
  void createOSSClient ();

private:
  boost::thread_specific_ptr<nsCOMPtr<ncIOSSClient>> _ossClientPtr;

  nsCOMPtr<ossGatewayInterface>   _ossGateway;
};

#endif // __NC_ACS_DATA_STORE_H__
