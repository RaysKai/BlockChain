package serialize


type SerializeStream interface{

}

type ISerialize interface{
	//Serialize/Deserialize
	Serialize()(SerializeStream)
	Deserialize(s SerializeStream)

	//
	ToString()(string)
}