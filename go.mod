module github.com/labulakalia/wazero_net

go 1.24.0

require (
	github.com/cloudsoda/go-smb2 v0.0.0-20250228001242-d4c70e6251cc
	github.com/jlaffaye/ftp v0.2.0
	github.com/tetratelabs/wazero v1.9.0
	golang.org/x/crypto v0.37.0
)

require github.com/pkg/sftp v1.13.9

require (
	github.com/cloudsoda/sddl v0.0.0-20250224235906-926454e91efc // indirect
	github.com/geoffgarside/ber v1.1.0 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/goidentity/v6 v6.0.1 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.4 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/kr/fs v0.1.0 // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
)

replace github.com/pkg/sftp => github.com/labulakalia/sftp v1.13.10-0.20250421063436-b983016d5069

replace github.com/cloudsoda/go-smb2 => github.com/labulakalia/go-smb2 v0.0.0-20250421065043-77cbbfc100dc
