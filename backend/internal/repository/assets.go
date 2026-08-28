package repository

type Image struct { ID, Name, Namespace, AccessLevel, Digest, Checksum, UpdatedAt string }
type Object struct { Name, Bucket, Prefix, Size, StorageType, UpdatedAt string }

func Images(projectID, namespace, name string, page int) ([]Image, int, bool) {
	if projectID == "" { return nil, 0, false }; all := []Image{{"image-1","training-base","magiclab","private","sha256:demo-image-v1","sha256:demo-image-v1","2026-08-28T00:00:00Z"},{"image-2","lpy-q1-gym","locomotion","private","sha256:gym-v1","sha256:gym-v1","2026-08-27T06:30:18Z"}}; out:=all[:0]; for _, v:=range all { if namespace!=""&&v.Namespace!=namespace {continue}; if name!=""&&v.Name!=name {continue}; out=append(out,v) }; return out,len(out),true
}

func Objects(projectID, bucket, prefix string) ([]Object, int, bool) { if projectID=="" {return nil,0,false}; items:=[]Object{{"admin",bucket,prefix,"-","-","-"},{"public",bucket,prefix,"-","-","-"}}; return items,len(items),true }
