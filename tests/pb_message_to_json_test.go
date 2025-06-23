package main

import (
	"github.com/eiiches/mysql-protobuf-functions/internal/dedent"
	"github.com/eiiches/mysql-protobuf-functions/internal/testutils"
	"testing"
)

func TestMessageToJson(t *testing.T) {
	p := testutils.NewProtoTestSupport(t, map[string]string{
		"main.proto": dedent.Pipe(`
			|syntax = "proto3";
			|message TestMess {
			|    int32 int32_field = 1;
			|}
		`),
	})

	AssertThatCall(t, "pb_descriptor_set_load(?, ?)", "a", p.GetSerializedFileDescriptorSet()).ShouldSucceed()
	defer func() {
		AssertThatCall(t, "pb_descriptor_set_delete(?)", "a").ShouldSucceed()
	}()

	data := p.JsonToProtobuf(".TestMess", dedent.Pipe(`
		|{"int32Field": 123}
	`))

	RunTestThatExpression(t, "pb_message_to_json(?, ?, ?)", "a", ".TestMess", data).IsEqualToJson(dedent.Pipe(`
		|{"int32Field": 123}
	`))
}
