数据库表描述

1.t_acs_custom_perm
该表用以存储gns路径的权限配置信息。
文档的创建和eacp的权限请求会产生权限配置记录。
evfsacs会使用这些记录来检查gns对象的操作合法性。

2.t_acs_doc
该表用以存储objectId（用户id，部门id，群组id）与gns路径的关系。
主要用来转化gns路径的显示名称。

3.t_acs_owner
该表用来存储gns路径与所有者的关系。
用在判断用户是否是某个gns路径的所有者，比如用户做出权限配置动作，打开回收站操作，编辑群组操作

4.t_acs_group_share
该表用来记录群组的基本信息，比如名称
