package hba

type Options struct {
  str string
  kv  map[string]string
}

func nextoption(str string) (key string, value string, rest string) {

}

func NewOptionsFromString(str string) (o Options, err error){
  o.str = str
  str = strings.Trim(str)
  o.kv = make(map[string]string)
  for {
    
  }
  return o, nil
}
