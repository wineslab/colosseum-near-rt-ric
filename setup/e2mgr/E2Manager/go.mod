module e2mgr

require (
	gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common v1.0.35
	gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities v1.0.35
	gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader v1.0.35
	gerrit.o-ran-sc.org/r/ric-plt/sdlgo v0.5.2
	github.com/asaskevich/govalidator v0.0.0-20190424111038-f61b66f89f4a // indirect
	github.com/go-ozzo/ozzo-validation v3.5.0+incompatible
	github.com/golang/protobuf v1.3.4
	github.com/gorilla/mux v1.7.0
	github.com/magiconair/properties v1.8.1
	github.com/pelletier/go-toml v1.5.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/viper v1.4.0
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/stretchr/testify v1.4.0
	go.uber.org/multierr v1.2.0 // indirect
	go.uber.org/zap v1.11.0
	golang.org/x/net v0.0.0-20191021144547-ec77196f6094 // indirect
	golang.org/x/sys v0.0.0-20191105231009-c1f44814a5cd // indirect
	google.golang.org/appengine v1.6.1 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/apimachinery v0.17.0
	k8s.io/client-go v0.17.0
)

replace gerrit.o-ran-sc.org/r/ric-plt/sdlgo => gerrit.o-ran-sc.org/r/ric-plt/sdlgo.git v0.5.2
