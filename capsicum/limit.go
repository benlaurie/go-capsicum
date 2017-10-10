package capsicum

import "os"

func LimitStd() error {
	r, err := CapRightsInit(CAP_READ)
	if err != nil {
		return err
	}
	err = CapRightsLimit(os.Stdin, r)
	if err != nil {
		return err
	}

	r, err = CapRightsInit(CAP_WRITE)
	if err != nil {
		return err
	}
	err = CapRightsLimit(os.Stdout, r)
	if err != nil {
		return err
	}
	err = CapRightsLimit(os.Stderr, r)
	if err != nil {
		return err
	}

	return nil
}
