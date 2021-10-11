package consts

var MsgFlags = map[int]string {
	AUTH_STUDENT: "user",
	AUTH_TEACHER: "teacher",
	AUTH_ADMIN:   "admin",

	SUCCESS : "ok",
	ERROR : "fail",
	INVALID_PARAMS : "请求参数错误",
	ACCESS_DENIED : "当前权限不够",
	ERROR_EXIST_USER : "已存在该用户",
	ERROR_NOT_EXIST_USER : "该用户不存在",
	ERROR_AUTH_CHECK_TOKEN_FAIL : "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT : "Token已超时",
	ERROR_AUTH_TOKEN : "Token生成失败",
	ERROR_AUTH : "用户名或密码错误",
	ERROR_AUTH_PERMISSION_DENIED : "用户权限不足",
	ERROR_DEPLOY_CREATE: "deployment创建失败",
	ERROR_DEPLOY_UPDATE: "deployment更新失败",
	ERROR_SERVICE_CREATE: "service创建失败",
	ERROR_SERVICE_UPDATE: "service更新失败",
	ERROR_YAML_READ: "读取yaml文件失败",
	ERROR_YAML_CONVERT: "转换yaml文件失败",
	ERROR_YAML_UNMARSHAL: "yaml数据编出失败",
	ERROR_DEPLOY_ALREADY_EXIST : "deployment已存在",
	ERROR_DEPLOY_NOT_EXIST: "deployment不存在",
	ERROR_SERVICE_ALREADY_EXIST: "service已存在",
	ERROR_SERVICE_NOT_EXIST: "service不存在",
	ERROR_POD_NOT_EXIST: "pod不存在",
	ERROR_RESOURCE_OUT: "用户剩余资源不足",
	ERROR_RESOURCEQUOTA_GET: "用户资源配置获取失败",
	ERROR_RESOURCEQUOTA_UPDATE :"用户资源配置更新失败",
	ERROR_DB_EXIST : "数据库已存在",
	ERROR_DB_EXEC : "创建数据库脚本执行失败，请检查数据库名称或密码复杂度",
	ERROR_DB_DELETE : "数据库删除失败，请刷新检查",
	ERROR_DB_NOT_EXIST : "数据库不存在",
	ERROR_SSH_DECODE : "ssh解码失败",
	ERROR_DELETE_NS : "删除命名空间失败",
	ERROR_FILE_CREATE : "服务器创建文件失败",
	ERROR_FILE_WRITE : "服务器写入文件失败",
	ERROR_FILE_COPY : "文件传进容器失败",
	ERROR_FILE_SHELL_COPY : "脚本传进容器失败",
	ERROR_FILE_SHELL_RUN : "脚本运行失败,请保证压缩包是由单个文件夹压缩得到",
	ERROR_PVC_CREATE : "PVC创建失败",
	ERROR_PVC_EXIST : "该PVC已存在",
	ERROR_PVC_NOT_EXIST : "PVC不存在",
	ERROR_USER_CREATE : "创建用户失败",
	ERROR_IMAGE_PUSH : "添加镜像失败",
	ERROR_CLASS_USER_NOT_EXIST: "部分学生不在数据库中，请检查学生是否注册",
	ERROR_CLASS_END: "课程已结束，请重新设置课程到期时间",
	ERROR_CLASS_TASK_EXIST: "任务已存在，请分配不同任务",
	ERROR_CLASS_TASK_SUBMIT: "任务作业提交失败",
	ERROR_CLASS_EXIST: "课程英文名已存在，请重命名",
	ERROR_CLASS_TASK_END: "任务已结束，无法提交",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}