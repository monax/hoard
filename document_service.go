package hoard

import (
	"io"

	"github.com/monax/hoard/v6/api"
	"github.com/monax/hoard/v6/grant"
	"github.com/monax/hoard/v6/meta"
)

func (service *hoardService) Download(grt *grant.Grant, srv api.Document_DownloadServer) error {
	doc, salt, err := GetDocument(service.gs, grt)
	if err != nil {
		return err
	}

	return SendDocument(srv, doc, salt, service.cs)
}

func (service *hoardService) Upload(srv api.Document_UploadServer) error {
	doc, spec, salt, err := ReceiveDocumentAndGrant(srv)
	if err != nil {
		return err
	}

	grt, err := PostDocument(service.gs, doc, spec, salt)
	if err != nil {
		return err
	}

	return srv.SendAndClose(grt)
}

type DocumentSender interface {
	Send(*api.PlaintextAndMeta) error
}

func SendDocument(srv DocumentSender, doc *meta.Document, salt []byte, cs int) error {
	out := new(api.PlaintextAndMeta)
	out.Input = &api.PlaintextAndMeta_Meta{Meta: doc.Meta}
	if err := srv.Send(out); err != nil {
		return err
	}

	out.Input = &api.PlaintextAndMeta_Plaintext{
		Plaintext: &api.Plaintext{
			Input: &api.Plaintext_Salt{
				Salt: salt,
			},
		},
	}
	if err := srv.Send(out); err != nil {
		return err
	}

	data := doc.GetData()
	for i := 0; i < len(data); i += cs {
		if i+cs > len(data) {
			out.Input = &api.PlaintextAndMeta_Plaintext{
				Plaintext: &api.Plaintext{
					Input: &api.Plaintext_Data{
						Data: data[i:len(data)],
					},
				},
			}
		} else {
			out.Input = &api.PlaintextAndMeta_Plaintext{
				Plaintext: &api.Plaintext{
					Input: &api.Plaintext_Data{
						Data: data[i : i+cs],
					},
				},
			}
		}
		if err := srv.Send(out); err != nil {
			return err
		}
	}

	return nil
}

type DocumentAndGrantSender interface {
	Send(*api.PlaintextAndGrantSpecAndMeta) error
}

func SendDocumentAndGrant(srv DocumentAndGrantSender, doc *meta.Document, salt []byte, spec *grant.Spec, cs int) error {
	out := new(api.PlaintextAndGrantSpecAndMeta)
	out.Input = &api.PlaintextAndGrantSpecAndMeta_Meta{Meta: doc.Meta}
	if err := srv.Send(out); err != nil {
		return err
	}

	out.Input = &api.PlaintextAndGrantSpecAndMeta_PlaintextAndGrantSpec{
		PlaintextAndGrantSpec: &api.PlaintextAndGrantSpec{
			Input: &api.PlaintextAndGrantSpec_GrantSpec{
				GrantSpec: spec,
			},
		},
	}
	if err := srv.Send(out); err != nil {
		return err
	}

	out.Input = &api.PlaintextAndGrantSpecAndMeta_PlaintextAndGrantSpec{
		PlaintextAndGrantSpec: &api.PlaintextAndGrantSpec{
			Input: &api.PlaintextAndGrantSpec_Plaintext{
				Plaintext: &api.Plaintext{
					Input: &api.Plaintext_Salt{
						Salt: salt,
					},
				},
			},
		},
	}
	if err := srv.Send(out); err != nil {
		return err
	}

	data := doc.GetData()
	for i := 0; i < len(data); i += cs {
		if i+cs > len(data) {
			out.Input = &api.PlaintextAndGrantSpecAndMeta_PlaintextAndGrantSpec{
				PlaintextAndGrantSpec: &api.PlaintextAndGrantSpec{
					Input: &api.PlaintextAndGrantSpec_Plaintext{
						Plaintext: &api.Plaintext{
							Input: &api.Plaintext_Data{
								Data: data[i:len(data)],
							},
						},
					},
				},
			}
		} else {
			out.Input = &api.PlaintextAndGrantSpecAndMeta_PlaintextAndGrantSpec{
				PlaintextAndGrantSpec: &api.PlaintextAndGrantSpec{
					Input: &api.PlaintextAndGrantSpec_Plaintext{
						Plaintext: &api.Plaintext{
							Input: &api.Plaintext_Data{
								Data: data[i : i+cs],
							},
						},
					},
				},
			}
		}
		if err := srv.Send(out); err != nil {
			return err
		}
	}

	return nil
}

type DocumentAndGrantReceiver interface {
	Recv() (*api.PlaintextAndGrantSpecAndMeta, error)
}

func ReceiveDocumentAndGrant(srv DocumentAndGrantReceiver) (*meta.Document, *grant.Spec, []byte, error) {
	doc := new(meta.Document)
	spec := new(grant.Spec)
	var salt []byte
	for {
		d, err := srv.Recv()
		if err != nil {
			if err == io.EOF {
				return doc, spec, salt, nil
			}

			return nil, nil, nil, err
		}

		switch x := d.GetInput().(type) {
		case *api.PlaintextAndGrantSpecAndMeta_Meta:
			doc.Meta = x.Meta
		case *api.PlaintextAndGrantSpecAndMeta_PlaintextAndGrantSpec:
			switch y := x.PlaintextAndGrantSpec.GetInput().(type) {
			case *api.PlaintextAndGrantSpec_GrantSpec:
				spec = y.GrantSpec
			case *api.PlaintextAndGrantSpec_Plaintext:
				switch z := y.Plaintext.GetInput().(type) {
				case *api.Plaintext_Salt:
					salt = z.Salt
				case *api.Plaintext_Data:
					doc.Data = append(doc.Data, z.Data...)
				}
			}
		}
	}
}

type DocumentReceiver interface {
	Recv() (*api.PlaintextAndMeta, error)
}

func ReceiveDocument(srv DocumentReceiver) (*meta.Document, []byte, error) {
	doc := new(meta.Document)
	var salt []byte
	for {
		d, err := srv.Recv()
		if err != nil {
			if err == io.EOF {
				return doc, salt, nil
			}

			return nil, nil, err
		}

		switch x := d.GetInput().(type) {
		case *api.PlaintextAndMeta_Meta:
			doc.Meta = x.Meta
		case *api.PlaintextAndMeta_Plaintext:
			switch y := x.Plaintext.GetInput().(type) {
			case *api.Plaintext_Salt:
				salt = y.Salt
			case *api.Plaintext_Data:
				doc.Data = append(doc.Data, y.Data...)
			}
		}
	}
}
