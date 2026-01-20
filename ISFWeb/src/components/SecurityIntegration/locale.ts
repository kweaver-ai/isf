import i18n from '@/core/i18n';

export default i18n([
    [
        '初始化配置',
        '初始化設定',
        'Initialization',
    ],
    [
        '确定',
        '確定',
        'OK',
    ],
    [
        '取消',
        '取消',
        'Cancel',
    ],
    [
        '用户密级策略',
        '用戶密級原則',
        'Secret Level Policy',
    ],
    [
        '非密',
        '非密',
        'Unclassified',
    ],
    [
        '内部',
        '內部',
        'Internal',
    ],
    [
        '秘密',
        '秘密',
        'Secret',
    ],
    [
        '机密',
        '機密',
        'Confidential',
    ],
    [
        '绝密',
        '絕密',
        'Top Secret',
    ],
    [
        '公开',
        '公開',
        'Public',
    ],
    [
        '密级列表：',
        '密級列表：',
        'Secret Level List',
    ],
    [
        '设置',
        '設定',
        'Set',
    ],
    [
        '在保存初始化配置之前，请先完成“自定义密级”的设置',
        '在儲存初始化配置之前，請先完成"自訂密級"的設置',
        "Before saving the initial configuration, please set 'Custom Secret Level' first.",
    ],
    [
        '密码策略',
        '密碼原則',
        'Password Policy',
    ],
    [
        '密码有效期：',
        '密碼有效期間：',
        'Expiration',
    ],
    [
        '天',
        '天',
        'day(s)',
    ],
    [
        '个月',
        '個月',
        'month(s)',
    ],
    [
        '永久有效',
        '永久有效',
        'Permanent',
    ],
    [
        '密码仅在指定时间段内有效，若超过该有效期，则需要修改密码，否则无法登录',
        '密碼僅在指定時間段內有效，若超過該有效期，則需要變更密碼，否則無法登入',
        'Password is only valid within the specified period. You will need to change the password after it expires or you can not log in.',
    ],
    [
        '密码强度：',
        '密碼強度：',
        'Strength',
    ],
    [
        '强密码',
        '強式密碼',
        'Strong',
    ],
    [
        '弱密码',
        '弱式密碼',
        'Weak',
    ],
    [
        '强密码格式：密码长度至少为',
        '強式密碼格式:密碼長度至少為',
        'Strong Passwords Format: at least ',
    ],
    [
        '个字符，需同时包含 大小写英文字母、数字与特殊字符',
        '個字元，需同時包含 大小寫英文字母、數字與特殊字元',
        'characters, and must contain numbers, letters in both uppercase and lowercase, and special characters',
    ],
    [
        '强密码长度至少为${num}个字符，请重新输入。',
        '強式密碼長度至少為${num}個字元，請重新輸入。',
        'Strong Passwords should be at least ${num} characters. Please re-enter.',
    ],
    [
        '弱密码格式：密码长度至少为6个字符',
        '弱式密碼格式：密碼長度至少為6個字元',
        'Weak Passwords Format: at least 6 characters',
    ],
    [
        '启用密码错误锁定：',
        '啟用密碼錯誤鎖定：',
        'Enalbe Password Error Lockdown: ',
    ],
    [
        '用户在任一情况下，密码连续输错',
        '使用者在任一情况下，密码连续输错',
        'An account will be locked after',
    ],
    [
        '次，则账号将被锁定',
        '次，則帳戶將會被鎖定',
        ' times of failed password attempts.',
    ],
    [
        '账号被锁定后，',
        '帳戶被鎖定后，',
        'The lock of an account will be unlocked in',
    ],
    [
        '分钟后自动解锁或由管理员解锁',
        '分鐘後自動解鎖或由管理員解鎖',
        ' minutes or Admin will unlock it.',
    ],
    [
        '请从低到高依次定义新的密级：',
        '請從低到高依次定義新的密級：',
        'Please define Secret Level from low to high:',
    ],
    [
        '新增密级',
        '新增密級',
        'Add Secret Level',
    ],
    [
        '密级',
        '密級',
        'Secret Level',
    ],
    [
        '操作',
        '操作',
        'Operation',
    ],
    [
        '初始化后将无法更改已设置的用户密级，请确认你的操作。',
        '初始化後將無法變更已設定的用戶密級，請確認您的操作。',
        'Once initialized, the Secret Levels cannot be changed. Please confirm your operation.',
    ],
    [
        '密码错误次数范围为1~${num}，请重新输入。',
        '密碼錯誤次數範圍為1~${num}，請重新輸入。',
        'Password attempts should be within 1~${num}, please re-enter.',
    ],
    [
        '账号锁定时间范围为10~180分钟，请重新输入。',
        '帳戶鎖定时间范围为10~180分钟，请重新输入。',
        'The time range should be within 10~180 minutes, please re-enter.',
    ],
    [
        '不能包含 / : * ? " < > | 特殊字符，请重新输入。',
        '不能包含 / : * ? " < > | 特殊字元，請重新輸入。',
        'Cannot contain special characters as / : * ? " < > | , please re-enter.',
    ],
    [
        '该密级名称已存在。',
        '該密級名稱已存在。',
        'The same name already exists.',
    ],
    [
        '您最多只能设置11个密级。',
        '您最多隻能設定11個密級。',
        'At most 11 Secret Levels can be added.',
    ],
    [
        '将文件密级从低到高自定义为“${secu}”成功',
        '將文件密級從低到高自訂為“${secu}”成功',
        'Customize file Secret Levels from lowest to highest "${secu}" successfully',
    ],
    [
        '将用户密级从低到高自定义为“${secu}”成功',
        '將用戶密級從低到高自訂為“${secu}”成功',
        'Customize user Secret Levels from lowest to highest "${secu}" successfully',
    ],
    [
        '将系统密级设置为 “${level}” 成功',
        '將系統密級設定為 “${level}” 成功',
        'Set system Secret Level as “${level}” successfully',

    ],
    [
        '原密级为“${originalLevel}”',
        '原密級為“${originalLevel}”',
        'Original Secret Level is “${originalLevel}”.',
    ],
    [
        '设置 密码策略 成功',
        '設定 密碼原則 成功',
        'Set password policy successfully',
    ],
    [
        '启用 密码错误锁定，最大连续输错密码次数为${passwdErrCnt}次，自动解锁时间为${passwdLockTime}分钟',
        '啟用 密碼錯誤鎖定，最大連續輸錯密碼次數為${passwdErrCnt}次，自动解锁时间为${passwdLockTime}分钟',
        'Enable password error lock. Maximum continous password attempts: ${passwdErrCnt} times. Lock period: ${passwdLockTime} minutes.',
    ],
    [
        '关闭 密码错误锁定，最大连续输错密码次数为${passwdErrCnt}次，自动解锁时间为${passwdLockTime}分钟',
        '關閉 密碼錯誤鎖定，最大連續輸錯密碼次數為${passwdErrCnt}次，自动解锁时间为${passwdLockTime}分钟',
        'Close password error lock. Maximum continous password attempts: ${passwdErrCnt} times. Lock period: ${passwdLockTime} minutes.',
    ],
    [
        '强密码，密码长度至少为${passwdCnt}个字符',
        '強式密碼，密碼長度至少為${passwdCnt}個字元',
        'Strong Passwords. Its length should be at least ${passwdCnt} characters.',
    ],
    [
        '密级名不能为空。',
        '密級名不能為空。',
        'Secret Level name is required.',
    ],
    [
        '因受父域登录策略管控，无法执行此操作。请在父域的【多文域管理-策略同步】页面关闭密码强度设置才可正常使用。',
        '因受父網域登入原則管控，無法執行此操作。請在父網域的【多文域管理-原則同步】頁面關閉密碼強度設定才可正常使用。',
        'We can\'t do that under such parent domain control. Please turn off [Password Strength] setting on [Doc Domains-Policy Sync] first.',
    ],
    [
        '取消设置 禁止同一账号多地同时登录 成功，（仅对使用Windows客户端登录进行限制，使用其他客户端登录不受限制）',
        '取消設定 禁止同一帳戶多地同時登入 成功，（僅對使用Windows用戶端登入進行限制，使用其它終端登入不受限制）',
        'Cancel Setting. Prohit simultaneous login with the same account from multiple locations successfully (Only available for Windows Client).',
    ],
    [
        '禁止同一帐号多地同时登录（仅对使用Windows客户端登录进行限制，使用其它终端登录不受限制）',
        '禁止同一帳戶多地同時登入（僅對使用Windows用戶端登入進行限制，使用其它終端登入不受限制）',
        'Prohit simultaneous login with the same account from multiple locations successfully (Only available for Windows Client).',
    ],
    [
        '系统保护等级',
        '系統保護等級',
        'System Protection Levels',
    ],
    [
        '系统保护等级：',
        '系統保護等級：',
        'Protection Level:',
    ],
    [
        '秘密级',
        '秘密級',
        'Secret'
    ],
    [
        '机密级一般',
        '機密級一般',
        'Confidential Standard'
    ],
    [
        '机密级增强',
        '機密級增強',
        'Confidential Enhanced'
    ],
    [
        '登录策略',
        '登入原則',
        'Login Policy',
    ],
])