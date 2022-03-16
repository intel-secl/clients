module intel/isecl/lib/clients/v4

require (
	github.com/google/uuid v1.1.2
	github.com/gorilla/mux v1.7.3
	github.com/sirupsen/logrus v1.4.0
	github.com/stretchr/testify v1.2.2
	intel/isecl/lib/common/v4 v4.2.0-Beta
)

replace intel/isecl/lib/common/v4 => github.com/intel-secl/common/v4 v4.2.0-Beta
