package version

import (
	"strconv"

	"github.com/goccy/go-json"
	"gitlab.galaxy123.cloud/base/public/helper/ipx"
	"gitlab.galaxy123.cloud/base/public/helper/strx"
)

type Version struct {
	Server    string `json:"server"`     // 服务名称
	BuildTime string `json:"build_time"` // 编译时间
	CommitId  string `json:"commit_id"`  // 提交gitID
	Branch    string `json:"branch"`     // 代码分支
	Listen    string `json:"listen"`     // 监听的端口和地址
}

func NewVersion(server, buildTime, commitId, branch string) *Version {

	return &Version{
		Server:    server,
		BuildTime: buildTime,
		CommitId:  commitId,
		Branch:    branch,
		Listen:    ipx.GetClientIp() + ":",
	}
}

func (r *Version) WithPort(port int) {
	r.Listen += strconv.Itoa(port)
}

func (r *Version) ToStr() string {
	marshal, _ := json.Marshal(r)
	return strx.B2s(marshal)
}

func (r *Version) ToBytes() []byte {
	marshal, _ := json.Marshal(r)
	return marshal
}
