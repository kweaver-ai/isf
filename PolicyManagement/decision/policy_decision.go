package decision

import (
	"context"
	_ "embed"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"policy_mgnt/dependency"
	"policy_mgnt/general"
	"policy_mgnt/infra/outbox"
	"policy_mgnt/network"
	"policy_mgnt/utils/models"

	"policy_mgnt/utils/gocommon/api"

	cerrors "policy_mgnt/utils/gocommon/v2/errors"
	clog "policy_mgnt/utils/gocommon/v2/log"
	outboxer "policy_mgnt/utils/gocommon/v2/outbox"
	cutils "policy_mgnt/utils/gocommon/v2/utils"

	jsoniter "github.com/json-iterator/go"
	"github.com/mitchellh/mapstructure"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/storage"
	"github.com/open-policy-agent/opa/storage/inmem"
)

type PolicyDecision interface {
	IncrementalUpdatePolicy(ctx context.Context, userId string, departmentPaths []string) error
	InitOpaData() error
	NetWorkDecision(ctx context.Context, reqBody models.Accessor) (effect bool, err error)
	ClientSignDecision(ctx context.Context, clientType string) (effect bool, err error)
	PublishChange(userId string, departmentPaths []string) error
	PublishInit() error

	// ------for debug------
	ReadOPAData(ctx context.Context, path string) (interface{}, error)
}

//todo 新增策略决策接口

var (
	policyOnce      sync.Once
	policySingleton *policyDecision
)

type policyDecision struct {
	abstractDriven dependency.AbstractDriven
	logger         clog.Logger
	opaStore       storage.Store
	decisionQuery  rego.PreparedEvalQuery
	queryQuery     rego.PreparedEvalQuery
	outbox         *outboxer.Outboxer
	mgnt           network.Management
	gmgnt          general.Management
}

// NewDocLibAccessPolicy 初始化管理实例
func NewPolicyDecision() PolicyDecision {
	mgnt, _ := network.NewManagement()
	gmgnt, _ := general.NewManagement()
	policyOnce.Do(func() {
		policySingleton = &policyDecision{
			abstractDriven: dependency.NewAbstractDriven(),
			logger:         clog.NewLogger(),
			outbox:         outbox.NewOutBoxer(),
			mgnt:           mgnt,
			gmgnt:          gmgnt,
		}
		policySingleton.init()
	})
	return policySingleton
}

var (
	//go:embed network/allow.rego
	NetworkAllowRego string
)

var (
	PolicyInitDataChannel   = "policy_management.policy_data.init"
	PolicyUpdateDataChannel = "policy_management.policy_data.update"
)

func (d *policyDecision) init() {
	data, err := d.getOPAData()
	if err != nil {
		d.logger.Panic(err)
	}
	if err := d.initOPA(data); err != nil {
		d.logger.Panic(err)
	}

	d.compareVersionTimer()
}

func (p *policyDecision) InitOpaData() error {
	data, err := p.getOPAData()
	if err != nil {
		p.logger.Panic(err)
	}
	if err := p.initOPA(data); err != nil {
		p.logger.Panic(err)
	}
	return err
}

func (d *policyDecision) initOPA(data map[string]any) error {
	store := inmem.NewFromObject(data)

	decisionQuery, err := rego.New(
		rego.Query("res = data.sign_in_policy.client_restriction"),
		rego.Store(store),
		rego.Strict(true),
	).PrepareForEval(context.Background())
	if err != nil {
		return err
	}

	queryCompiler := ast.MustCompileModules(map[string]string{
		"allow.rego": NetworkAllowRego,
	})
	queryQuery, err := rego.New(
		rego.Query("res = data.network.accessible"),
		rego.Compiler(queryCompiler),
		rego.Store(store),
		rego.Strict(true),
	).PrepareForEval(context.Background())
	if err != nil {
		return err
	}

	d.opaStore = store
	d.decisionQuery = decisionQuery
	d.queryQuery = queryQuery

	return nil
}

func toJSONThenToInterface(data any) (any, error) {
	byteData, err := jsoniter.Marshal(data)
	if err != nil {
		return nil, err
	}

	var result any
	if err := jsoniter.Unmarshal(byteData, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *policyDecision) compareVersionTimer() {
	go func() {
		interval := 5 * time.Minute
		timer := time.NewTicker(interval)
		for {
			<-timer.C
			for {
				storeData, err := p.ReadOPAData(context.Background(), "/")
				if err != nil {
					continue
				}

				orderedStoreData, err := toJSONThenToInterface(storeData)
				if err != nil {
					continue
				}

				dbData, err := p.getOPAData()
				if err != nil {
					continue
				}

				orderedDBData, err := toJSONThenToInterface(dbData)
				if err != nil {
					continue
				}

				if !reflect.DeepEqual(orderedStoreData, orderedDBData) {
					p.logger.Infoln("detected difference, start to sync policy opa data from database to memory")
					if err := p.initOPA(dbData); err != nil {
						p.logger.Errorln(err, "retry in ten seconds")
						time.Sleep(time.Second * 10)
						continue
					}
					p.logger.Infoln("policy opa data sync ended successfully")
				}
				break
			}
		}
	}()
}

func (p *policyDecision) getNetworkRestrictionState() (result bool, err error) {
	res, _, err := p.gmgnt.ListPolicy(0, 1, []string{"network_restriction"})
	if err != nil {
		return
	}

	value := res[0].Value
	var networkRestriction general.NetworkResitriction
	err = json.Unmarshal(value, &networkRestriction)
	if err != nil {
		return
	}
	result = networkRestriction.IsEnabled
	return
}

// 获取未设置网段策略访问者开关状态
func (p *policyDecision) getNoNetworkPolicyAccessorState() (result bool, err error) {
	res, _, err := p.gmgnt.ListPolicy(0, 1, []string{"no_network_policy_accessor"})
	if err != nil {
		return
	}

	value := res[0].Value
	var noNetworkPolicyAccessor general.NoNetworkPolicyAccessor
	err = json.Unmarshal(value, &noNetworkPolicyAccessor)
	if err != nil {
		return
	}
	result = noNetworkPolicyAccessor.IsEnabled
	return
}

// 获取需要更新到opa上的数据
func (p *policyDecision) getOPAData() (map[string]interface{}, error) {
	// 1、获取数据
	result := make(map[string]interface{})
	is_enabled, err := p.getNetworkRestrictionState()
	if err != nil {
		return result, err
	}

	netWorkData := make(map[string]interface{})
	// 如果开关开启，查询具体数据
	// 如果开关关闭，不查询具体数据
	if is_enabled {
		if err := p.mgnt.GetNetworkData(netWorkData); err != nil {
			return result, err
		}
		no_policy_enabled, err := p.getNoNetworkPolicyAccessorState()
		if err != nil {
			return result, err
		}
		netWorkData["no_policy_accessor_enabled"] = no_policy_enabled
	}
	netWorkData["is_enabled"] = is_enabled

	// 获取客户端登录选项数据
	clientResData := make(map[string]interface{})
	var value general.ClientRestriction
	var resValue models.ClientRestrictionConfig
	client_res, _, err := p.gmgnt.ListPolicy(0, 1, []string{(&general.ClientRestriction{}).Name()})
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(client_res[0].Value, &value)
	resValue.PcWeb = value.PCWEB
	resValue.MobileWeb = value.MobileWEB
	resValue.Windows = value.Windows
	resValue.Mac = value.Mac
	resValue.Android = value.Android
	resValue.IOS = value.IOS
	resValue.Linux = value.Linux
	clientResData[(&general.ClientRestriction{}).Name()] = resValue
	result["network_info"] = netWorkData
	result["sign_in_policy"] = clientResData

	return result, err
}

func (p *policyDecision) ReadOPAData(ctx context.Context, path string) (interface{}, error) {
	ph, ok := storage.ParsePathEscaped(path)
	if !ok {
		return nil, cerrors.ErrBadRequestPublic(&cerrors.ErrorInfo{Cause: "invalid path" + path})
	}

	txn, err := p.opaStore.NewTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err == nil {
			if err = p.opaStore.Commit(ctx, txn); err != nil {
				p.logger.Errorf("opa store commit error: %v", err)
			}
		} else {
			p.opaStore.Abort(ctx, txn)
			p.logger.Errorf("opa store txn aborted due to %v", err)
		}
	}()

	data, err := p.opaStore.Read(ctx, txn, ph)
	if err != nil {
		return nil, err
	}
	return data, nil
}

const (
	Ipv4Type       = "ipv4"
	Ipv6Type       = "ipv6"
	opaNetworkInfo = "/network_info"
)

// 增量更新用户所属部门
func (p *policyDecision) updateUserDepartmentRelation(userId string, departmentPaths []string) (content []models.OPADataPatch, err error) {
	db, err := api.ConnectDB()
	if err != nil {
		return
	}
	var updateInfo bool
	for _, v := range departmentPaths {
		res := strings.Split(v, "/")
		var netData models.NetworkAccessorRelation
		db.Select("f_id").Where("f_accessor_id in (?)", res).First(&netData)
		if netData.ID != 0 {
			updateInfo = true
		}
	}
	var relations []models.OPAUpdateInfo
	if updateInfo {
		var relation models.OPAUpdateInfo
		rows, err := db.Raw(`select
								a.f_accessor_id,
								b.f_ip_type
							from
								t_network_accessor_relation as a
							left join t_network_restriction as b on
								a.f_network_id = b.f_id
							where
								a.f_accessor_type = 'department'`).Rows()
		if err != nil {
			return content, fmt.Errorf("select data info error: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&relation.AccessorId, &relation.IpType)
			if err != nil {
				return content, fmt.Errorf("rows Scan error: %w", err)
			}
			relations = append(relations, relation)
		}
	}
	// 查看当前用户是否本身已经存在策略
	var userNetData models.NetworkAccessorRelation
	db.Select("f_id").Where("f_accessor_id = ?", userId).First(&userNetData)
	if len(relations) != 0 {
		// 调用Usermanagement接口获取所属部门
		departIds, err := p.abstractDriven.GetBelongDepartByUserId(userId)
		if err != nil {
			return content, fmt.Errorf("GetBelongDepartByUserId: %w", err)
		}
		var ipv4Ids []string
		var ipv6Ids []string
		for _, v := range relations {
			for _, value := range departIds {
				if v.IpType == "ipv4" && v.AccessorId == value {
					ipv4Ids = append(ipv4Ids, v.AccessorId)
				} else if v.IpType == "ipv6" && v.AccessorId == value {
					ipv6Ids = append(ipv6Ids, v.AccessorId)
				}
			}
		}
		// 增量添加用户到opa上
		var ipv4_content []models.OPADataPatch
		var ipv6_content []models.OPADataPatch
		if len(ipv4Ids) != 0 {
			ipv4_content = []models.OPADataPatch{{PatchOP: storage.AddOp, Path: "/network_info/ipv4/users/" + userId + "/departments", Value: ipv4Ids}}
			if userNetData.ID == 0 {
				ipv4_content = []models.OPADataPatch{{PatchOP: storage.AddOp, Path: "/network_info/ipv4/users/" + userId, Value: map[string]interface{}{"departments": ipv4Ids}}}
			}

		} else if len(ipv4Ids) == 0 && len(ipv6Ids) == 0 {
			ipv4_content = []models.OPADataPatch{{PatchOP: storage.RemoveOp, Path: "/network_info/ipv4/users/" + userId + "/departments", Value: map[string]interface{}{}}}
			if userNetData.ID == 0 {
				ipv4_content = []models.OPADataPatch{{PatchOP: storage.RemoveOp, Path: "/network_info/ipv4/users/" + userId, Value: map[string]interface{}{}}}
			}
		}
		if ipv4_content != nil {
			content = append(content, ipv4_content...)
		}
		if len(ipv6Ids) != 0 {
			ipv6_content = []models.OPADataPatch{{PatchOP: storage.AddOp, Path: "/network_info/ipv6/users/" + userId + "/departments", Value: ipv6Ids}}
			if userNetData.ID == 0 {
				ipv6_content = []models.OPADataPatch{{PatchOP: storage.AddOp, Path: "/network_info/ipv6/users/" + userId, Value: map[string]interface{}{"departments": ipv6Ids}}}
			}
		} else if len(ipv4Ids) == 0 && len(ipv6Ids) == 0 {
			ipv6_content = []models.OPADataPatch{{PatchOP: storage.RemoveOp, Path: "/network_info/ipv6/users/" + userId + "/departments", Value: map[string]interface{}{}}}
			if userNetData.ID == 0 {
				ipv6_content = []models.OPADataPatch{{PatchOP: storage.RemoveOp, Path: "/network_info/ipv6/users/" + userId, Value: map[string]interface{}{}}}
			}
		}
		if ipv6_content != nil {
			content = append(content, ipv6_content...)
		}

	}
	return
}

func (p *policyDecision) incrementalUpdateData(ctx context.Context, data []models.OPADataPatch) (err error) {
	txn, err := p.opaStore.NewTransaction(ctx, storage.WriteParams)
	if err != nil {
		return fmt.Errorf("NewTransaction: %w", err)
	}

	defer func() {
		if err == nil {
			if err = p.opaStore.Commit(ctx, txn); err != nil {
				p.logger.Errorf("opa store commit error: %v", err)
			}
		} else {
			p.opaStore.Abort(ctx, txn)
			p.logger.Errorf("opa store txn aborted due to %v", err)
		}
	}()

	for i := range data {
		path, ok := storage.ParsePathEscaped(data[i].Path)
		if !ok {
			return fmt.Errorf("ParsePathEscaped: %w", stderrors.New(cutils.ContactStr("invalid path", data[i].Path)))
		}
		// 所以此处忽略掉这个错误
		if err = p.opaStore.Write(ctx, txn, data[i].PatchOP, path, data[i].Value); err != nil && !storage.IsNotFound(err) {
			return fmt.Errorf("Write: %w", err)
		}
	}

	return nil
}

// 发送消息
func (p *policyDecision) PublishChange(userId string, departmentPaths []string) error {
	var updateData models.UserDepartRelation
	updateData.UserId = userId
	updateData.DepartmentPaths = departmentPaths
	v, err := jsoniter.Marshal(updateData)
	if err != nil {
		return fmt.Errorf("jsoniter.Marshal: %w", err)
	}

	// FIXME 同一个事务
	if err := p.outbox.Send(&outboxer.OutboxMessage{
		Payload: v,
		Options: outboxer.DynamicValues{outbox.Channel: PolicyUpdateDataChannel},
	}); err != nil {
		return fmt.Errorf("d.outbox.Send: %w", err)
	}

	p.logger.Infoln(PolicyUpdateDataChannel, userId, departmentPaths)
	return nil
}

func (p *policyDecision) PublishInit() error {
	// FIXME 同一个事务
	if err := p.outbox.Send(&outboxer.OutboxMessage{
		Payload: []byte(""), // 不允许为空, 所以填一个空字符串
		Options: outboxer.DynamicValues{outbox.Channel: PolicyInitDataChannel},
	}); err != nil {
		return fmt.Errorf("d.outbox.Send: %w", err)
	}

	p.logger.Infoln(PolicyInitDataChannel)
	return nil
}

func (p *policyDecision) IncrementalUpdatePolicy(ctx context.Context, userId string, departmentPaths []string) error {
	// 增量更新
	cont, err := p.updateUserDepartmentRelation(userId, departmentPaths)
	if err != nil {
		return fmt.Errorf("calculateIncrementalUpdateData: %w", err)
	}

	if err = p.incrementalUpdateData(ctx, cont); err != nil {
		return fmt.Errorf("incrementalUpdateData: %w", err)
	}

	return nil
}

func (p *policyDecision) NetWorkDecision(ctx context.Context, decisionOPAObjs models.Accessor) (res bool, err error) {
	results, err := p.queryQuery.Eval(ctx, rego.EvalInput(decisionOPAObjs))
	if err != nil {
		p.logger.Errorln(err)
		return
	} else if len(results) == 0 {
		return false, fmt.Errorf("NetWorkDecision: %w", stderrors.New("opa return empty"))
	}
	result := results[0].Bindings["res"]
	p.logger.Debugf("result: %#v", result)
	if err = mapstructure.Decode(result, &res); err != nil {
		p.logger.Errorln(err)
		return
	}

	return
}

func (p *policyDecision) ClientSignDecision(ctx context.Context, clientType string) (res bool, err error) {
	obj := make(map[string]bool)
	results, err := p.decisionQuery.Eval(ctx)
	if err != nil {
		p.logger.Errorln(err)
		return
	} else if len(results) == 0 {
		return false, fmt.Errorf("ClientSignDecision: %w", stderrors.New("opa return empty"))
	}

	result := results[0].Bindings["res"]
	p.logger.Debugf("result: %#v", result)
	if err = mapstructure.Decode(result, &obj); err != nil {
		p.logger.Errorln(err)
		return
	}
	if _, ok := obj[clientType]; !ok {
		return false, cerrors.ErrBadRequestPublic(&cerrors.ErrorInfo{Cause: "invalid client_type:" + clientType})
	}
	return obj[clientType], err
}
