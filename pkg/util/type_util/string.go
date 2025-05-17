package type_util

type TString string

func (receiver TString) String() string {
	return string(receiver)
}

func (receiver TString) Bytes() []byte {
	return []byte(receiver)
}

func (receiver TString) IsEmpty() bool {
	return 0 == len(receiver)
}

func (receiver TString) Equals(other TString) bool {
	return receiver == other
}
