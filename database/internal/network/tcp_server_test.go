package network

/*func TestTCPServer(t *testing.T) {
	request := []byte("request")
	response := []byte("response")

	tests := map[string]struct {
		client func(string)
		result int
	}{
		"single query": {
			client: func(address string) {
				connection, err := net.Dial("tcp", address)
				require.NoError(t, err)

				defer func(connection net.Conn) {
					err := connection.Close()
					require.NoError(t, err)
				}(connection)

				_, err = connection.Write(request)
				require.NoError(t, err)

				buffer := make([]byte, 0, 2048)
				_, err = connection.Read(buffer)
				require.NoError(t, err)

				require.True(t, reflect.DeepEqual(response, buffer))
			},
		},
		"multiple queries": {},
	}
}
*/
